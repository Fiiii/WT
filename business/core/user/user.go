package user

import (
	"context"
	"fmt"
	"github.com/Fiiii/WT/business/repository/store/user"
	"go.uber.org/zap"
	"time"
)

// Core manages the set of API's for user access.
type Core struct {
	userStore user.Store
}

//NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger) Core {
	return Core{
		userStore: user.NewStore(log),
	}
}

// Create inserts a new user into the database.
func (c Core) Create(ctx context.Context, nu user.NewUser, now time.Time) (*user.User, error) {
	usr, err := c.userStore.Create(ctx, nu, now)
	if err != nil {
		return &user.User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

// Update replaces a user document in the database.
func (c Core) Update(ctx context.Context, userID string, uu user.UpdateUser, now time.Time) error {
	if err := c.userStore.Update(ctx, userID, uu, now); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}
	return nil
}

// Delete removes a user from the database.
func (c Core) Delete(ctx context.Context, userID string) error {
	if err := c.userStore.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

// List retrieves a list of existing users from the database.
func (c Core) List(ctx context.Context) ([]*user.User, error) {
	users, err := c.userStore.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return users, nil
}

// QueryByID gets the specified user from the database.
func (c Core) QueryByID(ctx context.Context, userID string) (*user.User, error) {
	usr, err := c.userStore.Query(ctx, userID)
	if err != nil {
		return &user.User{}, fmt.Errorf("query: %w", err)
	}
	return usr, nil
}
