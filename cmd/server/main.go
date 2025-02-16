package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"log/slog"

	configs "github.com/danicatalao/notifier/configs/server"
	"github.com/danicatalao/notifier/internal/forecast"
	"github.com/danicatalao/notifier/internal/scheduled_notification"
	"github.com/danicatalao/notifier/internal/user"
	postgres "github.com/danicatalao/notifier/pkg/database"
	"github.com/gin-gonic/gin"
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
	}
	fmt.Printf("%+v\n", cfg)

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
	forecastHandler := forecast.NewForecastHandler(forecastApiClient, forecastService, log)

	userRepository := user.NewUserRepository(db, log)
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

	r.Run(net.JoinHostPort("", cfg.Http.Port))
}
