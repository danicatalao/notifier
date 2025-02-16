package forecast

import (
	"context"
	"fmt"
	"strconv"

	l "github.com/danicatalao/notifier/internal/logger"
)

type ForecastService interface {
	GetForecastAndWave(ctx context.Context, cityName string) (*ForecastWave, error)
}

type forecast_service struct {
	provider ForecastApiClient
	log      l.Logger
}

func NewForecastService(p ForecastApiClient, l l.Logger) *forecast_service {
	return &forecast_service{provider: p, log: l}
}

func (s *forecast_service) GetForecastAndWave(ctx context.Context, cityName string) (*ForecastWave, error) {
	city, err := s.provider.GetCity(ctx, cityName)
	if err != nil {
		return nil, fmt.Errorf("error trying to get city name: %w", err)
	}
	cityId := strconv.Itoa(city.Id)

	s.log.InfoContext(ctx, "Getting weather and wave forecast", "city", cityName)

	forecast, err := s.provider.GetCityForecast(ctx, cityId)
	if err != nil {
		return nil, fmt.Errorf("error trying to get city forecast: %w", err)
	}

	wave, err := s.provider.GetWaveForecast(ctx, cityId, "0")
	if err != nil {
		return nil, fmt.Errorf("error trying to get city wave forecast: %w", err)
	}

	var waveReturn *WaveForecast
	forecastReturn := forecast

	if wave.Name != "undefined" {
		waveReturn = wave
	}

	return &ForecastWave{Forecast: forecastReturn, Wave: waveReturn}, nil
}
