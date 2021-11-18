package usersGrp

import (
	"context"
	"fmt"
	"net/http"

	userCore "github.com/Fiiii/WT/business/core/user"
	"github.com/Fiiii/WT/business/repository/store/user"
	"github.com/Fiiii/WT/foundation/web"
)

type Handlers struct {
	User userCore.Core
}

// List returns a list of users with paging.
func (h Handlers) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	users, err := h.User.List(ctx)
	if err != nil {
		return fmt.Errorf("unable to query for users: %w", err)
	}

	return web.Respond(ctx, w, users, http.StatusOK)
}

// QueryByID returns a user by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	usr, err := h.User.QueryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ID[%s]: %w", id, err)
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}

// Create adds a new user to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var nu user.NewUser
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	usr, err := h.User.Create(ctx, nu, v.Now)
	if err != nil {
		return fmt.Errorf("user[%+v]: %w", &usr, err)
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

// Update updates a user in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var upd user.UpdateUser
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	id := web.Param(r, "id")
	if err := h.User.Update(ctx, id, upd, v.Now); err != nil {
		return fmt.Errorf("ID[%s] User[%+v]: %w", id, &upd, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a user from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	if err := h.User.Delete(ctx, id); err != nil {
		return fmt.Errorf("ID[%s]: %w", id, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
