package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

// Handler definition for HTTP requests handling.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint for the application.
type App struct {
	shutdown chan os.Signal
	mux      *httptreemux.ContextMux
}

// NewApp created new application instance.
func NewApp(shutdown chan os.Signal) *App {
	mux := httptreemux.NewContextMux()
	return &App{
		shutdown: shutdown,
		mux:      mux,
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
	a.mux.ServeHTTP(w, r)
}

// Handle sets a handler function for a given HTTP method, path with wrapping by middleware (TODO)
func (a *App) Handle(method, group, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		// Pull the context
		ctx := r.Context()
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
	a.mux.Handle(method, finalPath, h)
}
