package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/Fiiii/WT/app/services/handlers"
	"github.com/Fiiii/WT/business/sys/storage/dynamodb"
	"github.com/Fiiii/WT/foundation/logger"
	"github.com/ardanlabs/conf/v2"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	// Construct the application logger.
	log, err := logger.New("WT-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
		}
		Database struct {
			Profile string `conf:"dev-profile"`
			Project string `conf:"default:WT"`
			Stage   string `conf:"dev"`
			Region  string `conf:"eu-central-1"`
		}
	}{
		Version: conf.Version{
			Build: "dev",
			Desc:  "",
		},
	}

	const prefix = "WT"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// =========================================================================
	// Database Support

	// Create connectivity to the database.
	log.Infow("startup", "status", "initializing database support")
	db, err := dynamodb.NewClient(
		cfg.Database.Project,
		cfg.Database.Stage,
		cfg.Database.Region,
		cfg.Database.Profile,
	)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	// =========================================================================
	// Start Debug Service

	log.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	// The Debug function returns a mux to listen and serve on for all the debug
	// related endpoints. This include the standard library endpoints.

	// Construct the mux for the debug calls.
	debugMux := handlers.DebugMux(build, log, db.Client)

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)

	// Signal to relay incoming signals
	//signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMuxConf := handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		DB:       db.Client,
	}

	apiMux := handlers.APIMux(apiMuxConf)

	// Construct a server to service the requests against the mux.
	httpServer := http.Server{
		Addr:    ":3000",
		Handler: apiMux,
	}

	// Channel for listening errors from http server
	serverErrors := make(chan error, 1)

	// Starting the server
	go func() {
		serverErrors <- httpServer.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	if err := <-serverErrors; err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
