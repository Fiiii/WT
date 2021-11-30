package web

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
)

// Handler definition for HTTP requests handling.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint for the application.
type App struct {
	shutdown   chan os.Signal
	ContextMux *httptreemux.ContextMux
	mw         []Middleware
}

// NewApp created new application instance.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	mux := httptreemux.NewContextMux()
	return &App{
		shutdown:   shutdown,
		ContextMux: mux,
		mw:         mw,
	}
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// ServeHTTP implements the http.Handler interface. In the future it will be used as
// entrypoint with core for telemetry metrics.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.ContextMux.ServeHTTP(w, r)
}

// Handle sets a handler function for a given HTTP method, path with wrapping by middleware
func (a *App) Handle(method, group, path string, handler Handler, mw ...Middleware) {

	// First wrap handler specific middleware.
	handler = wrapMiddleware(mw, handler)

	// Secondly wrap by application's general middleware.
	handler = wrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {
		// Pull the context
		ctx := r.Context()

		// Capture the parent request span from the context.
		//span := trace.SpanFromContext(ctx)

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}

		ctx = context.WithValue(ctx, key, &v)

		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	// Creates end path based on provided group
	finalPath := path
	if group != "" {
		finalPath = fmt.Sprintf("/%s%s", group, path)
	}

	// Final handle by using httptreemux
	a.ContextMux.Handle(method, finalPath, h)
}
