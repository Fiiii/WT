package handlers

import (
	calculatorapi "github.com/Fiiii/WT/app/services/handlers/calculator-api"
	"github.com/Fiiii/WT/foundation/web"
	"net/http"
	"os"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) http.Handler {
	app := web.NewApp(cfg.Shutdown)
	v1(app)
	return app
}

// v1 aggregates all routes to the single version.
func v1(app *web.App) {
	const version = "v1"

	chg := calculatorapi.Handler{}
	app.Handle(http.MethodGet, version, "/calculate/:a/:b", chg.Add)
}