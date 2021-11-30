package middleware

import (
	"context"
	"github.com/Fiiii/WT/foundation/web"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Logger - returns nested (onion) handler function by the middleware wrapper.
func Logger (log *zap.SugaredLogger) web.Middleware {
	m := func (handler web.Handler) web.Handler {
		h := func (ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}

			log.Infow("request started", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			// Call the next handler.
			err = handler(ctx, w, r)

			log.Infow("request completed", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))

			return err
		}

		return h
	}
	return m
}


