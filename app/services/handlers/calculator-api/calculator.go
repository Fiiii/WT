package calculator_api

import (
	"context"
	"github.com/Fiiii/WT/foundation/web"
	"net/http"
	"strconv"
)

type Handler struct {

}

func (h *Handler) Add(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	a, _ := strconv.Atoi(web.Param(r, "a"))
	b, _ := strconv.Atoi(web.Param(r, "b"))
	return web.Respond(ctx, w, a+b, http.StatusOK)
}

