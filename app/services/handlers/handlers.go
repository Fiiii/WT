package handlers

import (
	"github.com/Fiiii/WT/app/services/handlers/v1/users-api"
	"github.com/Fiiii/WT/business/core/user"
	"net/http"
	"os"

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
	app := web.NewApp(cfg.Shutdown)

	// Load routes with previously initiated configuration.
	v1(app, cfg)
	return app
}

// v1 aggregates all routes to the single version.
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	// Register user management and authentication endpoints.
	ugh := users_api.Handlers{
		User: user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodGet, version, "/users", ugh.List)
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID)
	app.Handle(http.MethodPost, version, "/users", ugh.Create)
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update)
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete)
}
