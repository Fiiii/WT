package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/Fiiii/WT/app/services/wt-api/handlers"
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
		log.Errorw("startup", "ERROR", err)
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
	// Configuration | options: ie. mask, noprint

	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s,mask"`
		}
		Database struct {
			Profile string `conf:"default:dev-profile"`
			Project string `conf:"default:WT"`
			Stage   string `conf:"default:dev"`
			Region  string `conf:"default:eu-central-1"`
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
	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	strConf, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("config string error: %w", err)
	}

	log.Infow("startup", "config", strConf)
	expvar.NewString("build").Set(build)

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
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := httpServer.Shutdown(ctx); err != nil {
			httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
