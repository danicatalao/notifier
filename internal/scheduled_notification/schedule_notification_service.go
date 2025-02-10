package scheduled_notification

import (
	"context"
	"fmt"
)

type ScheduledNotificationService interface {
	CreateScheduledNotification(ctx context.Context, sn *ScheduledNotification) error
}

type schedule_notification_service struct {
	repository ScheduledNotificationRepository
}

func NewScheduledNotificationService(r ScheduledNotificationRepository) *schedule_notification_service {
	return &schedule_notification_service{repository: r}
}

func (s *schedule_notification_service) CreateScheduledNotification(ctx context.Context, sn *ScheduledNotification) error {
	err := s.repository.Create(ctx, sn)
	if err != nil {
		return fmt.Errorf("error trying to schedule a notification: %w", err)
	}
	return nil
}
