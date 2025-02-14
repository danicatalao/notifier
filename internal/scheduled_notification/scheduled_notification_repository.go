package scheduled_notification

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	l "github.com/danicatalao/notifier/internal/logger"
	"github.com/danicatalao/notifier/pkg/database"
)

type ScheduledNotificationRepository interface {
	Create(ctx context.Context, schedule *ScheduledNotification) error
	GetDueNotifications(ctx context.Context, batchSize uint64) ([]ScheduledNotification, error)
	UpdateNotificationStatusWithTx(ctx context.Context, id int64, status NotificationStatus) error
}

type scheduled_notification_repository struct {
	db  *database.Service
	log l.Logger
}

func NewScheduledNotificationRepository(db *database.Service, l l.Logger) *scheduled_notification_repository {
	return &scheduled_notification_repository{db: db, log: l}
}

func (r *scheduled_notification_repository) Create(ctx context.Context, sn *ScheduledNotification) error {
	query, args, err := r.db.Builder.
		Insert(SCHEDULED_NOTIFICATION_TABLE).
		Columns("date, city_name, user_id, notification_type").
		Values(sn.Date, sn.CityName, sn.UserId, sn.NotificationType).
		ToSql()

	if err != nil {
		return fmt.Errorf("could not build query: %w", err)
	}
	r.log.Info(ctx, "Executing sql statement", "sql", query, "args", args)

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}
	return nil
}

func (r *scheduled_notification_repository) GetDueNotifications(ctx context.Context, batchSize uint64) ([]ScheduledNotification, error) {
	var notifications []ScheduledNotification
	query, args, err := r.db.Builder.Select(
		"id",
		"status",
		"date",
		"city_name",
		"user_id",
		"notification_type",
		"created_at",
		"updated_at",
	).
		From(SCHEDULED_NOTIFICATION_TABLE).
		Where(squirrel.Eq{"status": StatusPending}).
		Where(squirrel.LtOrEq{"date": time.Now()}).
		OrderBy("date ASC").
		Limit(batchSize).
		Suffix("FOR UPDATE SKIP LOCKED").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("could not build query: %w", err)
	}
	r.log.Info(ctx, "Executing sql statement", "sql", query, "args", args)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notification ScheduledNotification
		err := rows.Scan(
			&notification.Id,
			&notification.Status,
			&notification.Date,
			&notification.CityName,
			&notification.UserId,
			&notification.NotificationType,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return notifications, nil
}

func (r *scheduled_notification_repository) UpdateNotificationStatusWithTx(ctx context.Context, id int64, status NotificationStatus) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query, args, err := r.db.Builder.
		Update(SCHEDULED_NOTIFICATION_TABLE).
		Set("status", status).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("could not build query: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update notification %d: %w", id, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
