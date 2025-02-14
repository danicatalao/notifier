package notification_producer

import (
	"context"
	"fmt"
	"log"
	"time"

	l "github.com/danicatalao/notifier/internal/logger"
	"github.com/danicatalao/notifier/internal/scheduled_notification"
	sn "github.com/danicatalao/notifier/internal/scheduled_notification"
	"github.com/danicatalao/notifier/pkg/rabbitmq"
)

type Worker struct {
	rabbitService rabbitmq.Service
	repository    sn.ScheduledNotificationRepository
	log           l.Logger
	pollInterval  time.Duration
	batchSize     uint64
}

type NotificationMessage struct {
	CityName string
	UserId   int64
}

func NewWorker(rs rabbitmq.Service, r sn.ScheduledNotificationRepository, p time.Duration, b uint64, l l.Logger) *Worker {
	return &Worker{
		rabbitService: rs,
		repository:    r,
		pollInterval:  p,
		batchSize:     b,
		log:           l,
	}
}

func (w *Worker) publishNotification(ctx context.Context, n sn.ScheduledNotification) error {
	routingKey := fmt.Sprintf("%s.notifications", n.NotificationType)

	msg := rabbitmq.Message{
		RoutingKey: routingKey,
		Body:       &NotificationMessage{n.CityName, n.UserId},
	}

	w.log.Info(ctx, "Queuing notification", "message", msg)
	return w.rabbitService.Publish(ctx, msg)
}

func (w *Worker) Start(ctx context.Context) error {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.processDueNotifications(ctx); err != nil {
				log.Printf("Error processing notifications: %v", err)
			}
		}
	}
}

func (w *Worker) processDueNotifications(ctx context.Context) error {
	notifications, err := w.repository.GetDueNotifications(ctx, w.batchSize)
	if err != nil {
		return fmt.Errorf("failed to fetch notifications: %v", err)
	}

	for _, notification := range notifications {
		if err := w.publishNotification(ctx, notification); err != nil {
			w.log.Error(ctx, "Failed to publish notification", "notification", notification.Id, "error", err.Error())
			if err := w.repository.UpdateNotificationStatusWithTx(ctx, notification.Id, scheduled_notification.StatusFailed); err != nil {
				w.log.Error(ctx, "Failed to update notification status: %v", err)
			}
			continue
		}

		if err := w.repository.UpdateNotificationStatusWithTx(ctx, notification.Id, scheduled_notification.StatusSent); err != nil {
			log.Printf("Failed to update notification status: %v", err)
		}
	}

	return nil
}
