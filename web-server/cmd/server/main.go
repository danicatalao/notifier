package main

import (
	"fmt"
	"net"
	"os"

	"github.com/danicatalao/notifier/web-server/configs"
	"github.com/danicatalao/notifier/web-server/internal/user"
	postgres "github.com/danicatalao/notifier/web-server/pkg/database"
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

	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	// HTTP Server
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		userHandler.AddUserRoutes(v1)
	}

	r.Run(net.JoinHostPort("", cfg.HTTP.Port))
}
