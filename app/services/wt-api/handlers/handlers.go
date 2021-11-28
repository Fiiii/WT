package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/Fiiii/WT/app/services/wt-api/handlers/debug/checkgrp"
	"github.com/Fiiii/WT/app/services/wt-api/handlers/v1/productsGrp"
	"github.com/Fiiii/WT/app/services/wt-api/handlers/v1/usersGrp"
	"github.com/Fiiii/WT/business/core/product"
	"github.com/Fiiii/WT/business/core/user"
	"github.com/Fiiii/WT/business/middleware"
	"github.com/Fiiii/WT/foundation/web"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	DB       *dynamodb.Client
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) http.Handler {
	app := web.NewApp(
		cfg.Shutdown,
		middleware.Logger(cfg.Log),
	)

	// Load routes with previously initiated configuration.
	v1(app, cfg)
	return app
}

// v1 aggregates all routes to the single version.
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	// Register user management endpoints.
	ugh := usersGrp.Handlers{
		User: user.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodGet, version, "/users", ugh.List)
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID)
	app.Handle(http.MethodPost, version, "/users", ugh.Create)
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update)
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete)

	// Register product management endpoints.
	pgh := productsGrp.Handlers{
		Product: product.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodGet, version, "/products", pgh.Query)
	app.Handle(http.MethodGet, version, "/products/:id", pgh.QueryByID)
	app.Handle(http.MethodPost, version, "/products", pgh.Create)
	app.Handle(http.MethodPut, version, "/products/:id", pgh.Update)
	app.Handle(http.MethodDelete, version, "/products/:id", pgh.Delete)
}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service. This bypassing the use of the
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger, db *dynamodb.Client) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
		DB:    db,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}
