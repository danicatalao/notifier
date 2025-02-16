package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/danicatalao/notifier/internal/forecast"
	l "github.com/danicatalao/notifier/internal/logger"
	"github.com/danicatalao/notifier/internal/user"
)

type HttpClient interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type webhook_service struct {
	userService     user.UserService
	forecastService forecast.ForecastService
	log             l.Logger
	httpClient      HttpClient
}

func NewWebhookService(u user.UserService, f forecast.ForecastService, l l.Logger, h HttpClient) *webhook_service {
	return &webhook_service{
		userService:     u,
		forecastService: f,
		log:             l,
		httpClient:      h,
	}
}

func (s *webhook_service) Send(ctx context.Context, id int64, cityName string) error {
	user, err := s.userService.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !user.Active {
		return fmt.Errorf("user is not accepting notifications")
	}

	if user.Webhook == nil || *user.Webhook == "" {
		return fmt.Errorf("user does not have a webhook")
	}
	forecastwave, err := s.forecastService.GetForecastAndWave(ctx, cityName)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(forecastwave)
	if err != nil {
		return fmt.Errorf("could not parse forecast response")
	}

	s.log.InfoContext(ctx, "Sending Webhook", "usuario", &user, "city", cityName, "content", string(jsonData))

	response, err := s.httpClient.Post(*user.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("could not request the webhook")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook request failed with code %d", response.StatusCode)
	}

	return nil
}
