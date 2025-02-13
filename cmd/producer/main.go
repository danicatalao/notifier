package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	configs "github.com/danicatalao/notifier/configs/producer"
	"github.com/danicatalao/notifier/internal/scheduled_notification"
	postgres "github.com/danicatalao/notifier/pkg/database"
	"github.com/danicatalao/notifier/pkg/rabbitmq"
)

func main() {
	ctx := context.Background()

	log := slog.Make(sloghuman.Sink(os.Stdout))

	// Loading .env variables into config
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatal(ctx, "Config error: %s", err)
	}
	fmt.Printf("%+v\n", cfg)

	messageBroker, err := rabbitmq.NewService(rabbitmq.Config{
		Url:            cfg.Rabbitmq.Url,
		ExchangeName:   cfg.Rabbitmq.ExchangeName,
		ReconnectDelay: cfg.Rabbitmq.ReconnectDelay * time.Second,
		MaxRetries:     cfg.Rabbitmq.MaxRetries,
	})
	defer messageBroker.Close()

	db, err := postgres.New(cfg.Pg.Url, cfg.Pg.ConnAttempts, cfg.Pg.ConnTimeoutMs)
	if err != nil {
		log.Fatal(ctx, "could not create connection pool on Postgres", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()
	log.Info(ctx, "Connection pool created on Postgres")

	scheduledNotificationRepository := scheduled_notification.NewScheduledNotificationRepository(db)
	scheduledNotificationService := scheduled_notification.NewScheduledNotificationService(scheduledNotificationRepository)

	fmt.Print("Hello World")

}
