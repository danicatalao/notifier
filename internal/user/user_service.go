package user

import (
	"context"
	"fmt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *AppUser) (int64, error)
	GetByID(ctx context.Context, id int64) (*AppUser, error)
}

type user_service struct {
	repository UserRepository
}

func NewUserService(r UserRepository) *user_service {
	return &user_service{repository: r}
}

func (s *user_service) CreateUser(ctx context.Context, user *AppUser) (int64, error) {
	id, err := s.repository.Create(ctx, user)
	if err != nil {
		return -1, fmt.Errorf("error trying to add a user: %w", err)
	}
	return id, nil
}

func (s *user_service) GetByID(ctx context.Context, id int64) (*AppUser, error) {
	user, err := s.repository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error trying to get a user by Id: %w", err)
	}
	return user, nil
}
