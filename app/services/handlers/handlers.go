package handlers

import (
	"github.com/Fiiii/WT/foundation/web"
	"net/http"
)

func APIMux() http.Handler {
	app := web.NewApp()
	return app
}
