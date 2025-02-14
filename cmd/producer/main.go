package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	configs "github.com/danicatalao/notifier/configs/producer"
	"github.com/danicatalao/notifier/internal/notification_producer"
	"github.com/danicatalao/notifier/internal/scheduled_notification"
	postgres "github.com/danicatalao/notifier/pkg/database"
	"github.com/danicatalao/notifier/pkg/rabbitmq"
)

func main() {
	ctx := context.Background()

	log := slog.Make(sloghuman.Sink(os.Stdout))

	// Loading .env variables into config
	cfg, err := configs.NewConfig(".env")
	if err != nil {
		log.Fatal(ctx, "Config error", "error", err)
	}
	fmt.Printf("%+v\n", cfg)

	messageBroker, err := rabbitmq.NewService(rabbitmq.Config{
		Url:            cfg.Rabbitmq.Url,
		ExchangeName:   cfg.Rabbitmq.ExchangeName,
		ReconnectDelay: cfg.Rabbitmq.ReconnectDelay * time.Second,
		MaxRetries:     cfg.Rabbitmq.MaxRetries,
	})
	if err != nil {
		log.Fatal(ctx, "failed to create RabbitMQ service", "error", err)
	}
	defer messageBroker.Close()

	db, err := postgres.New(cfg.Pg.Url, cfg.Pg.ConnAttempts, cfg.Pg.ConnTimeoutMs)
	if err != nil {
		log.Fatal(ctx, "could not create connection pool on Postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Info(ctx, "Connection pool created on Postgres")

	scheduledNotificationRepository := scheduled_notification.NewScheduledNotificationRepository(db)

	worker := notification_producer.NewWorker(messageBroker, scheduledNotificationRepository, cfg.PollInterval, cfg.BatchSize, log)

	if err := worker.Start(ctx); err != nil {
		log.Fatal(ctx, "failed to start producer worker", "error", err)
	}

	fmt.Print("Hello World")
}
