package scheduled_notification

import (
	"context"
	"fmt"

	"github.com/danicatalao/notifier/pkg/database"
)

type ScheduleNotificationRepository interface {
	Create(ctx context.Context, schedule *ScheduledNotification) error
	//Get(ctx context.Context)
}

type scheduled_notification_repository struct {
	*database.Service
}

func NewScheduleNotificationRepository(db *database.Service) *scheduled_notification_repository {
	return &scheduled_notification_repository{db}
}

func (r *scheduled_notification_repository) Create(ctx context.Context, sn *ScheduledNotification) error {
	query, args, err := r.Builder.
		Insert(SCHEDULED_NOTIFICATION_TABLE).
		Columns("status, date, city_name, user_id, notification_type").
		Values(sn.Status, sn.Date, sn.CityName, sn.UserID, sn.NotificationType).
		ToSql()

	fmt.Printf("%s\n", query)

	if err != nil {
		return fmt.Errorf("user repository - adduser - could not build query: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}
	return nil
}
