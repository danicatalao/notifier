package configs

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		Pg               `env-prefix:"PG_" env-required:"true"`
		Http             `env-prefix:"HTTP_" env-required:"true"`
		ForecastProvider `env-prefix:"FORECAST_PROVIDER_" env-required:"true"`
	}

	Pg struct {
		Url           string `env:"URL"`
		ConnAttempts  int    `env:"CONN_ATTEMPTS"`
		ConnTimeoutMs int    `env:"CONN_TIMEOUT_MS"`
	}

	Http struct {
		Port string `env:"PORT" env-default:"8080"`
	}

	ForecastProvider struct {
		Url string `env:"URL"`
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
