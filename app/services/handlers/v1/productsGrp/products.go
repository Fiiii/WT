package productsGrp

import (
	"context"
	"fmt"
	productCore "github.com/Fiiii/WT/business/core/product"
	"github.com/Fiiii/WT/business/repository/store/product"
	"github.com/Fiiii/WT/foundation/web"
	"net/http"
)

type Handlers struct {
	Product productCore.Core
}

// Query returns a list of products.
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	products, err := h.Product.List(ctx)
	if err != nil {
		return fmt.Errorf("unable to query for products: %w", err)
	}

	return web.Respond(ctx, w, products, http.StatusOK)
}

// QueryByID returns a product by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	prod, err := h.Product.QueryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ID[%s]: %w", id, err)
	}

	return web.Respond(ctx, w, prod, http.StatusOK)
}

// Create adds a new product to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	prod, err := h.Product.Create(ctx, np, v.Now)
	if err != nil {
		return fmt.Errorf("creating new product, np[%+v]: %w", np, err)
	}

	return web.Respond(ctx, w, prod, http.StatusCreated)
}

// Update updates a product in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var upd product.UpdateProduct
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	id := web.Param(r, "id")
	if err := h.Product.Update(ctx, id, upd, v.Now); err != nil {
		return fmt.Errorf("ID[%s] Product[%+v]: %w", id, &upd, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a product from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	if err := h.Product.Delete(ctx, id); err != nil {
		return fmt.Errorf("ID[%s]: %w", id, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
