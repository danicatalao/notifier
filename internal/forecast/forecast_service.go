package forecast

import (
	"fmt"
	"strconv"
)

type ForecastService interface {
	GetForecastAndWave(cityName string) (*ForecastWave, error)
}

type forecast_service struct {
	provider ForecastApiClient
}

func NewForecastService(p ForecastApiClient) *forecast_service {
	return &forecast_service{provider: p}
}

func (s *forecast_service) GetForecastAndWave(cityName string) (*ForecastWave, error) {
	city, err := s.provider.GetCity(cityName)
	if err != nil {
		return nil, fmt.Errorf("error trying to get city name: %w", err)
	}
	cityId := strconv.Itoa(city.ID)
	fmt.Printf("%+v\n", city)

	forecast, err := s.provider.GetCityForecast(cityId)
	if err != nil {
		return nil, fmt.Errorf("error trying to get city forecast: %w", err)
	}
	fmt.Printf("%+v\n", forecast)

	wave, err := s.provider.GetWaveForecast(cityId, "0")
	if err != nil {
		return nil, fmt.Errorf("error trying to get city wave forecast: %w", err)
	}
	fmt.Printf("%+v\n", wave)

	var waveReturn *WaveForecast
	forecastReturn := forecast

	if wave.Name != "undefined" {
		waveReturn = wave
	}

	fmt.Printf("%+v\n", waveReturn)
	fmt.Printf("%+v\n", forecastReturn)

	return &ForecastWave{Forecast: forecastReturn, Wave: waveReturn}, nil
}
