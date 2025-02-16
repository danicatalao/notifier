package notification

import (
	"context"

	"github.com/danicatalao/notifier/internal/forecast"
	l "github.com/danicatalao/notifier/internal/logger"
	"github.com/danicatalao/notifier/internal/user"
)

type webhook_service struct {
	userService     user.UserService
	forecastService forecast.ForecastService
	log             l.Logger
}

func NewWebhookService(u user.UserService, f forecast.ForecastService, l l.Logger) *webhook_service {
	return &webhook_service{
		userService:     u,
		forecastService: f,
		log:             l,
	}
}

func (s *webhook_service) Send(ctx context.Context) error {
	s.log.InfoContext(ctx, "Webhook sent")
	return nil
}
