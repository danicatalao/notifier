package scheduled_notification

import (
	"context"
	"fmt"

	l "github.com/danicatalao/notifier/internal/logger"
)

type ScheduledNotificationService interface {
	CreateScheduledNotification(ctx context.Context, sn *ScheduledNotification) error
}

type schedule_notification_service struct {
	repository ScheduledNotificationRepository
	log        l.Logger
}

func NewScheduledNotificationService(r ScheduledNotificationRepository, l l.Logger) *schedule_notification_service {
	return &schedule_notification_service{repository: r, log: l}
}

func (s *schedule_notification_service) CreateScheduledNotification(ctx context.Context, sn *ScheduledNotification) error {
	err := s.repository.Create(ctx, sn)
	if err != nil {
		return fmt.Errorf("error trying to schedule a notification: %w", err)
	}
	return nil
}
