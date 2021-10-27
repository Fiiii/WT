package product

import (
	"context"
	"fmt"
	"github.com/Fiiii/WT/business/repository/store/product"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
	"time"
)

// Core manages the set of API's for product access.
type Core struct {
	productStore product.Store
}

func NewCore(log *zap.SugaredLogger, db *dynamodb.Client) Core {
	return Core{
		productStore: product.NewStore(log, db),
	}
}

// Create inserts a new product into the database.
func (c Core) Create(ctx context.Context, np product.NewProduct, now time.Time) (*product.Product, error)  {
	pd, err := c.productStore.Create(ctx, np, now)
	if err != nil {
		return &product.Product{}, fmt.Errorf("create: %w", err)
	}

	return pd, err
}

// Update replaces a product document in the database.
func (c Core) Update(ctx context.Context, productID string, uu product.UpdateProduct, now time.Time) error {
	if err := c.productStore.Update(ctx, productID, uu, now); err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

// Delete deletes a product from the database.
func (c Core) Delete(ctx context.Context, productID string) error  {
	if err := c.productStore.Delete(ctx, productID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// List retrieves a list of existing products from the database.
func (c Core) List(ctx context.Context) ([]*product.Product, error) {
	products, err := c.productStore.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return products, nil
}

// QueryByID gets the specified product from the database.
func (c Core) QueryByID(ctx context.Context, userID string) (*product.Product, error) {
	pd, err := c.productStore.QueryByID(ctx, userID)
	if err != nil {
		return &product.Product{}, fmt.Errorf("query: %w", err)
	}
	return pd, nil
}

