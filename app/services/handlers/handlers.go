package handlers

import (
	calculatorapi "github.com/Fiiii/WT/app/services/handlers/calculator-api"
	"github.com/Fiiii/WT/foundation/web"
	"net/http"
)

func APIMux() http.Handler {
	app := web.NewApp()
	return app
}


func v1(app *web.App) {
	const version = "v1"

	chg := calculatorapi.Handler{}

	app.Handle(http.MethodGet, version, "/calculate", chg.Add())
}