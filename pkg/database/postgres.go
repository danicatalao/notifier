package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

func New(url string, connAttempts int, connTimeoutMs int) (*Service, error) {
	for i := 1; i <= connAttempts; i++ {
		pool, err := pgxpool.New(context.Background(), url)
		if err == nil {
			pg := &Service{Pool: pool}
			pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
			return pg, nil
		}
		log.Print(fmt.Errorf("Error creating connection pool: %w", err))
		log.Printf("Trying to connect to Postgres -- Attempts left: %d", connAttempts-i)

		time.Sleep(time.Millisecond * time.Duration(connTimeoutMs))
	}
	return nil, fmt.Errorf("Postgres Service - Could not establish connection to Postgres")
}

func (p *Service) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
