package notification_producer

import (
	"context"
	"fmt"
	"log"
	"time"

	l "github.com/danicatalao/notifier/internal/logger"
	"github.com/danicatalao/notifier/internal/scheduled_notification"
	sn "github.com/danicatalao/notifier/internal/scheduled_notification"
	"github.com/danicatalao/notifier/pkg/database"
	"github.com/danicatalao/notifier/pkg/rabbitmq"
)

type Worker struct {
	databaseService *database.Service
	rabbitService   rabbitmq.Service
	pollInterval    time.Duration
	batchSize       int
	log             l.Logger
}

type NotificationMessage struct {
	CityName string
	UserId   int64
}

func NewWorker(db *database.Service, rs rabbitmq.Service, p time.Duration, b int, l l.Logger) *Worker {
	return &Worker{
		databaseService: db,
		rabbitService:   rs,
		pollInterval:    p,
		batchSize:       b,
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
	notifications, err := w.fetchDueNotifications(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch notifications: %v", err)
	}

	for _, notification := range notifications {
		if err := w.publishNotification(ctx, notification); err != nil {
			w.log.Error(ctx, "Failed to publish notification", "notification", notification.ID, "error", err.Error())
			if err := w.updateNotificationStatus(ctx, notification.ID, scheduled_notification.StatusFailed); err != nil {
				w.log.Error(ctx, "Failed to update notification status: %v", err)
			}
			continue
		}

		if err := w.updateNotificationStatus(ctx, notification.ID, scheduled_notification.StatusSent); err != nil {
			log.Printf("Failed to update notification status: %v", err)
		}
	}

	return nil
}
