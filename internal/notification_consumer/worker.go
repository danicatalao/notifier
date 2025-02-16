package notification_consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/danicatalao/notifier/internal/forecast"
	l "github.com/danicatalao/notifier/internal/logger"
	"github.com/danicatalao/notifier/internal/notification"
	"github.com/danicatalao/notifier/internal/user"
	"github.com/danicatalao/notifier/pkg/rabbitmq"
)

type NotificationMessage struct {
	CityName string
	UserId   int64
}

type HttpClient interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Worker struct {
	rabbitService       rabbitmq.Service
	userService         user.UserService
	forecastService     forecast.ForecastService
	notificationService notification.NotificationService
	log                 l.Logger
	queueName           string
}

func NewWorker(r rabbitmq.Service, u user.UserService, f forecast.ForecastService, l l.Logger, queue string, h HttpClient) (*Worker, error) {
	var n notification.NotificationService

	if queue == "webhook.notifications" {
		n = notification.NewWebhookService(u, f, l, h)
	} else {
		return nil, fmt.Errorf("notification type %s not supported", queue)
	}

	return &Worker{
		rabbitService:       r,
		userService:         u,
		forecastService:     f,
		notificationService: n,
		log:                 l,
		queueName:           queue,
	}, nil
}

func (w *Worker) handleMessage(ctx context.Context, body []byte) error {
	var message NotificationMessage
	err := json.Unmarshal(body, &message)
	if err != nil {
		return err
	}

	if err := w.notificationService.Send(ctx, message.UserId, message.CityName); err != nil {
		return err
	}

	return nil
}

func (w *Worker) Start(ctx context.Context) error {
	w.log.InfoContext(ctx, "Starting worker", "queue", w.queueName)

	deliveries, err := w.rabbitService.Consume(ctx, w.queueName)
	if err != nil {
		return err
	}

	w.log.InfoContext(ctx, "Worker started successfully", "queue", w.queueName)

	// Process messages until context is canceled
	for {
		select {
		case <-ctx.Done():
			w.log.InfoContext(ctx, "Context canceled, worker stopping", "queue", w.queueName)
			return ctx.Err()

		case delivery, ok := <-deliveries:
			if !ok {
				w.log.InfoContext(ctx, "Delivery channel closed, worker stopping", "queue", w.queueName)
				return nil
			}

			// Process the message
			w.log.DebugContext(ctx, "Received message", "queue", w.queueName, "bodySize", len(delivery.Body))

			err := w.handleMessage(ctx, delivery.Body)

			if err != nil {
				w.log.ErrorContext(ctx, "Failed to process message",
					"queue", w.queueName,
					"error", err)

				// Nack the message and don't requeue as we've already tried processing it
				if err := delivery.Nack(false, false); err != nil {
					w.log.ErrorContext(ctx, "Failed to nack message", "error", err)
				}
			} else {
				// Acknowledge successful processing
				if err := delivery.Ack(false); err != nil {
					w.log.ErrorContext(ctx, "Failed to acknowledge message", "error", err)
				}
			}
		}
	}
}
