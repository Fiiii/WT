// Package product contains product related CRUD functionality.
package product

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/Fiiii/WT/business/sys/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Store manages the set of API's for product access.
type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

// NewStore constructs a product store for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// Create adds a Product to the database. It returns the created Product with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, np NewProduct, now time.Time) (*Product, error) {
	if err := validate.Check(np); err != nil {
		return &Product{}, fmt.Errorf("validating data: %w", err)
	}

	prd := Product{
		ID:          validate.GenerateID(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}

	s.log.Infof("Creating new Product: %+v\n", prd)
	return &prd, nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func (s Store) Update(ctx context.Context, productID string, up UpdateProduct, now time.Time) error {
	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	prd, err := s.QueryByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("updating product productID[%s]: %w", productID, err)
	}

	if up.Name != nil {
		prd.Name = *up.Name
	}
	if up.Cost != nil {
		prd.Cost = *up.Cost
	}
	if up.Quantity != nil {
		prd.Quantity = *up.Quantity
	}

	prd.DateUpdated = now

	s.log.Infow("Updating specific product")
	return nil
}

// Delete removes the product identified by a given ID.
func (s Store) Delete(ctx context.Context, productID string) error {
	s.log.Infow("Deleting specific product")
	return nil
}

// List gets all Products from the database.
func (s Store) List(ctx context.Context) ([]*Product, error) {
	var products []*Product
	s.log.Infow("Querying for all products")

	products = append(products, &Product{})
	return products, nil
}

// QueryByID finds the product identified by a given ID.
func (s Store) QueryByID(ctx context.Context, productID string) (*Product, error) {
	if err := validate.CheckID(productID); err != nil {
		return &Product{}, errors.Wrap(err, "query single user error")
	}

	s.log.Infof("Querying for user by specific ID: %s\n", productID)
	return &Product{}, nil
}
