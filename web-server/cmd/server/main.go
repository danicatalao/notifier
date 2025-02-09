package main

import (
	"fmt"
	"os"

	config "github.com/danicatalao/notifier/web-server/configs"
	postgres "github.com/danicatalao/notifier/web-server/pkg/database"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Level = logrus.TraceLevel
	log.Out = os.Stdout

	// Loading .env variables into config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	fmt.Printf("%+v\n", cfg)

	pg, err := postgres.New(cfg.PG.URL, cfg.PG.ConnAttempts, cfg.PG.ConnTimeoutMs)
	if err != nil {
		log.Fatal("could not create connection pool on Postgres: %w", err, os.Stderr)
		os.Exit(1)
	}
	defer pg.Close()
	log.Info("Connection pool created on Postgres")

}
