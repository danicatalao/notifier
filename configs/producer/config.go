package configs

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Pg       `env-prefix:"PG_" env-required:"true"`
		Rabbitmq `env-prefix:"RABBITMQ_" env-required:"true"`
		Worker   `env-prefix:"WORKER_" env-required:"true"`
	}

	Pg struct {
		Url           string `env:"URL"`
		ConnAttempts  int    `env:"CONN_ATTEMPTS"`
		ConnTimeoutMs int    `env:"CONN_TIMEOUT_MS"`
	}

	Rabbitmq struct {
		Url            string        `env:"URL"`
		ExchangeName   string        `env:"EXCHANGE_NAME"`
		MaxRetries     int           `env:"MAX_RETRIES"`
		ReconnectDelay time.Duration `env:"RECONNECT_DELAY"`
	}

	Worker struct {
		PollInterval time.Duration `env:"POLL_INTERVAL"`
		BatchSize    uint64        `env:"BATCH_SIZE"`
	}
)

func NewConfig(filepath string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
