package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Fiiii/WT/business/sys/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Store manages the set of data layer access.
type Store struct {
	log *zap.SugaredLogger
}

// NewStore constructs new user Store.
func NewStore(log *zap.SugaredLogger) Store {
	return Store{
		log: log,
	}
}

// Create creates a new User based on NewUser input data.
func (s Store) Create(ctx context.Context, nu NewUser, now time.Time) (*User, error) {
	if err := validate.Check(nu); err != nil {
		return &User{}, fmt.Errorf("validating data: %w", err)
	}

	usr := User{
		Name:  nu.Name,
		Email: nu.Email,
		Roles: nu.Roles,
	}

	s.log.Infof("Creating new User: %+v\n", usr)
	return &usr, nil
}

// Query returns single user based on the userID.
func (s Store) Query(ctx context.Context, userID string) (*User, error) {
	if err := validate.CheckID(userID); err != nil {
		return &User{}, errors.Wrap(err, "query single user error")
	}

	s.log.Infof("Querying for user by specific ID: %s\n", userID)
	return &User{}, nil
}

// List returns all users.
func (s Store) List(ctx context.Context) ([]*User, error) {
	var users []*User
	s.log.Infow("Querying for all users")
	return users, nil
}

// Update updates concrete user by given userID.
func (s Store) Update(ctx context.Context, userID string, uu UpdateUser, now time.Time) error {
	if err := validate.Check(uu); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	s.log.Infow("Querying for specific user for the update")
	usr := &User{}

	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}

	s.log.Infow("Updating specific user")
	return nil
}

// Delete deletes concrete user by given userID.
func (s Store) Delete(ctx context.Context, userID string) error {
	s.log.Infow("Deleting specific user")
	return nil
}
