package user

import (
	"context"
	"fmt"

	l "github.com/danicatalao/notifier/internal/logger"

	"github.com/Masterminds/squirrel"
	"github.com/danicatalao/notifier/pkg/database"
)

type UserRepository interface {
	Create(ctx context.Context, user *AppUser) (int64, error)
	GetById(ctx context.Context, id int64) (*AppUser, error)
}

type user_repository struct {
	db  *database.Service
	log l.Logger
}

func NewUserRepository(db *database.Service, l l.Logger) *user_repository {
	return &user_repository{db: db, log: l}
}

func (r *user_repository) Create(ctx context.Context, u *AppUser) (int64, error) {
	query, args, err := r.db.Builder.
		Insert(APP_USER_TABLE).
		Columns("name, email, phone_number, webhook").
		Values(u.Name, u.Email, u.PhoneNumber, u.Webhook).
		ToSql()
	query += " RETURNING id"

	r.log.DebugContext(ctx, "Executing sql statement", "sql", query, "args", args)

	if err != nil {
		return -1, fmt.Errorf("could not build query: %w", err)
	}

	var id int64
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("queryrow failed: %w", err)
	}
	return id, nil
}

func (r *user_repository) GetById(ctx context.Context, id int64) (*AppUser, error) {
	var user AppUser
	query, args, err := r.db.Builder.Select(
		"id",
		"name",
		"email",
		"phone_number",
		"webhook",
		"active",
		"opt_out_date",
		"created_at",
		"updated_at",
	).
		From(APP_USER_TABLE).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("could not build query: %w", err)
	}
	r.log.DebugContext(ctx, "Executing sql statement", "sql", query, "args", args)

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.PhoneNumber,
		&user.Webhook,
		&user.Active,
		&user.OptOutDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	fmt.Printf("Usuario retornado pela query: %+v\n", user)

	return &user, nil
}
