package configs

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Pg               `env-prefix:"PG_" env-required:"true"`
		ForecastProvider `env-prefix:"FORECAST_PROVIDER_" env-required:"true"`
		Rabbitmq         `env-prefix:"RABBITMQ_" env-required:"true"`
		Queue            `env-prefix:"QUEUE_" env-required:"true"`
	}

	Pg struct {
		Url           string `env:"URL"`
		ConnAttempts  int    `env:"CONN_ATTEMPTS"`
		ConnTimeoutMs int    `env:"CONN_TIMEOUT_MS"`
	}

	Rabbitmq struct {
		Url            string        `env:"URL"`
		ExchangeName   string        `env:"EXCHANGE_NAME"`
		ExchangeType   string        `env:"EXCHANGE_TYPE"`
		MaxRetries     int           `env:"MAX_RETRIES"`
		ReconnectDelay time.Duration `env:"RECONNECT_DELAY"`
	}

	ForecastProvider struct {
		Url string `env:"URL"`
	}

	Queue struct {
		Name string `env:"NAME"`
	}
)

// func NewConfig(filepath string) (*Config, error) {
// 	cfg := &Config{}
// 	if err := cleanenv.ReadConfig(filepath, cfg); err != nil {
// 		return nil, err
// 	}
// 	return cfg, nil
// }

func NewConfig(filepath string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
