package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	fmt.Printf("%+v\n", db)
	return &user_repository{db}
}

func (r *user_repository) Create(ctx context.Context, u *AppUser) error {
	query, args, err := r.Builder.
		Insert(APP_USER_TABLE).
		Columns("name, email, phone_number, webhook").
		Values(u.Name, u.Email, u.PhoneNumber, u.Webhook).
		ToSql()
	query += " RETURNING id"

	fmt.Printf("%s\n", query)

	if err != nil {
		return fmt.Errorf("user repository - adduser - could not build query: %w", err)
	}

	var id int64
	err = r.Pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return fmt.Errorf("queryrow failed: %w", err)
	}
	fmt.Printf("ID do usuario: %d", id)
	return nil
}

func NewNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func NewNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}
