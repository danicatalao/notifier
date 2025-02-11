package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/danicatalao/notifier/configs"
	"github.com/danicatalao/notifier/internal/forecast"
	"github.com/danicatalao/notifier/internal/scheduled_notification"
	"github.com/danicatalao/notifier/internal/user"
	postgres "github.com/danicatalao/notifier/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Level = logrus.TraceLevel
	log.Out = os.Stdout

	// Loading .env variables into config
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	fmt.Printf("%+v\n", cfg)

	db, err := postgres.New(cfg.PG.URL, cfg.PG.ConnAttempts, cfg.PG.ConnTimeoutMs)
	if err != nil {
		log.Fatal("could not create connection pool on Postgres: %w", err, os.Stderr)
		os.Exit(1)
	}
	defer db.Close()
	log.Info("Connection pool created on Postgres")

	httpClient := &http.Client{Timeout: 10 * time.Second}

	forecastApiClient := forecast.NewForecastApiClient(httpClient, cfg.FORECAST_PROVIDER.URL)
	forecastService := forecast.NewForecastService(forecastApiClient)
	forecastHandler := forecast.NewForecastHandler(forecastApiClient, forecastService)

	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	scheduledNotificationRepository := scheduled_notification.NewScheduledNotificationRepository(db)
	scheduledNotificationService := scheduled_notification.NewScheduledNotificationService(scheduledNotificationRepository)
	scheduledNotificationHandler := scheduled_notification.NewScheduledNotificationHandler(scheduledNotificationService)

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
