package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
)

// Handler definition for HTTP requests handling.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint for the application.
type App struct {
	shutdown chan os.Signal
	mux      *http.ServeMux
}

// NewApp created new application instance.
func NewApp() *App {
	return &App{}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (a *App) Handle() {

}
