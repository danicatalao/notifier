package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"log/slog"

	configs "github.com/danicatalao/notifier/configs/consumer"
	"github.com/danicatalao/notifier/internal/forecast"
	"github.com/danicatalao/notifier/internal/notification_consumer"
	"github.com/danicatalao/notifier/internal/user"
	postgres "github.com/danicatalao/notifier/pkg/database"
	"github.com/danicatalao/notifier/pkg/rabbitmq"
	"github.com/lmittmann/tint"
)

func main() {
	ctx := context.Background()

	log := slog.New(tint.NewHandler(os.Stderr, nil))
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))

	// Loading .env variables into config
	cfg, err := configs.NewConfig(".env")
	if err != nil {
		log.ErrorContext(ctx, "Config error", "error", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", cfg)

	messageBroker, err := rabbitmq.NewService(rabbitmq.Config{
		Url:            cfg.Rabbitmq.Url,
		ExchangeName:   cfg.Rabbitmq.ExchangeName,
		ExchangeType:   cfg.ExchangeType,
		ReconnectDelay: cfg.Rabbitmq.ReconnectDelay,
		MaxRetries:     cfg.Rabbitmq.MaxRetries,
	}, log)
	if err != nil {
		log.ErrorContext(ctx, "failed to create RabbitMQ service", "error", err)
		os.Exit(1)
	}
	defer messageBroker.Close()
	log.InfoContext(ctx, "Connection established with RabbitMQ")

	db, err := postgres.New(cfg.Pg.Url, cfg.Pg.ConnAttempts, cfg.Pg.ConnTimeoutMs)
	if err != nil {
		log.ErrorContext(ctx, "could not create connection pool on Postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	log.InfoContext(ctx, "Connection pool created on Postgres")

	httpClient := &http.Client{Timeout: 10 * time.Second}

	forecastApiClient := forecast.NewForecastApiClient(httpClient, cfg.ForecastProvider.Url, log)
	forecastService := forecast.NewForecastService(forecastApiClient, log)

	userRepository := user.NewUserRepository(db, log)
	userService := user.NewUserService(userRepository)

	worker, err := notification_consumer.NewWorker(messageBroker, userService, forecastService, log, cfg.Queue.Name, httpClient)
	if err != nil {
		log.ErrorContext(ctx, "Could not create worker", "error", err)
		os.Exit(1)
	}

	if err := worker.Start(ctx); err != nil {
		log.ErrorContext(ctx, "failed to start consumer worker", "error", err)
		os.Exit(1)
	}
}
