package user

import (
	"context"
	"fmt"

	"github.com/danicatalao/notifier/web-server/pkg/database"
)

type UserRepository interface {
	Create(ctx context.Context, user *AppUser) error
	// List(ctx context.Context) ([]AppUser, error)
	// GetByID(ctx context.Context, id int64) (*AppUser, error)
	// Update(ctx context.Context, user *AppUser) error
	// Delete(ctx context.Context, id int64) error
}

type user_repository struct {
	*database.Service
}

func NewUserRepository(db *database.Service) *user_repository {
	return &user_repository{db}
}

func (r *user_repository) Create(ctx context.Context, u *AppUser) error {
	query, args, err := r.Builder.
		Insert(APP_USER_TABLE).
		Columns("name, email, phone_number, webhook, opt_out_date").
		Values(u.Name, u.Email, u.GetPhoneNumber(), u.GetWebhook(), u.GetOptOutDate()).
		ToSql()

	if err != nil {
		return fmt.Errorf("user repository - adduser - could not build query: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("user repository - adduser - could not insert User: %w", err)
	}

	return nil
}
