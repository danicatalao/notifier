package user

import (
	"context"
	"fmt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *AppUser) error
	// List(ctx context.Context) ([]AppUser, error)
	// GetByID(ctx context.Context, id int64) (*AppUser, error)
	// Update(ctx context.Context, user *AppUser) error
	// Delete(ctx context.Context, id int64) error
}

type user_service struct {
	repository UserRepository
}

func NewUserService(r UserRepository) *user_service {
	return &user_service{repository: r}
}

func (u *user_service) CreateUser(ctx context.Context, user *AppUser) error {
	err := u.repository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error trying to add a user: %w", err)
	}
	return nil
}
