package configs

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		PG                `env-prefix:"PG_" env-required:"true"`
		HTTP              `env-prefix:"HTTP_" env-required:"true"`
		FORECAST_PROVIDER `env-prefix:"FORECAST_PROVIDER_" env-required:"true"`
	}

	PG struct {
		URL           string `env:"URL"`
		ConnAttempts  int    `env:"CONN_ATTEMPTS"`
		ConnTimeoutMs int    `env:"CONN_TIMEOUT_MS"`
	}

	HTTP struct {
		Port string `env:"PORT" env-default:"8080"`
	}

	FORECAST_PROVIDER struct {
		URL string `env:"URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
