package user

import (
	"context"
	"fmt"
	"github.com/Fiiii/WT/business/sys/validate"
	"go.uber.org/zap"
	"time"
)

type Store struct {
	log *zap.SugaredLogger
}

func NewStore(log *zap.SugaredLogger) Store {
	return Store{
		log: log,
	}
}

func (s Store) Create(ctx context.Context, nu NewUser, now time.Time) (User, error) {
	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	usr := User{
		Name: nu.Name,
		Email: nu.Email,
		Roles: nu.Roles,
	}

	s.log.Infof("Creating new User: %+v\n", usr)
	return usr, nil
}
