package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	configs "github.com/danicatalao/notifier/configs/server"
	"github.com/danicatalao/notifier/internal/forecast"
	"github.com/danicatalao/notifier/internal/scheduled_notification"
	"github.com/danicatalao/notifier/internal/user"
	postgres "github.com/danicatalao/notifier/pkg/database"
	"github.com/gin-gonic/gin"
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

	db, err := postgres.New(cfg.PG.Url, cfg.PG.ConnAttempts, cfg.PG.ConnTimeoutMs)
	if err != nil {
		log.Fatal(ctx, "could not create connection pool on Postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Info(ctx, "Connection pool created on Postgres")

	httpClient := &http.Client{Timeout: 10 * time.Second}

	forecastApiClient := forecast.NewForecastApiClient(httpClient, cfg.FORECAST_PROVIDER.Url, log)
	forecastService := forecast.NewForecastService(forecastApiClient, log)
	forecastHandler := forecast.NewForecastHandler(forecastApiClient, forecastService, log)

	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	scheduledNotificationRepository := scheduled_notification.NewScheduledNotificationRepository(db, log)
	scheduledNotificationService := scheduled_notification.NewScheduledNotificationService(scheduledNotificationRepository, log)
	scheduledNotificationHandler := scheduled_notification.NewScheduledNotificationHandler(scheduledNotificationService, log)

	// HTTP Server
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		forecastHandler.AddForecastRoutes(v1)
		userHandler.AddUserRoutes(v1)
		scheduledNotificationHandler.AddNotificationRoutes(v1)
	}

	r.Run(net.JoinHostPort("", cfg.HTTP.Port))
}
