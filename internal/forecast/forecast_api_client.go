package forecast

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	l "github.com/danicatalao/notifier/internal/logger"
	"golang.org/x/net/html/charset"
)

type HttpClient interface {
	Get(s string) (resp *http.Response, err error)
}

type ForecastApiClient struct {
	httpClient HttpClient
	baseUrl    string
	log        l.Logger
}

func NewForecastApiClient(c HttpClient, s string, l l.Logger) ForecastApiClient {
	return ForecastApiClient{
		httpClient: c,
		baseUrl:    s,
		log:        l,
	}
}

func (c *ForecastApiClient) GetCity(ctx context.Context, name string) (*City, error) {
	url := fmt.Sprintf("%s/listaCidades?city=%s", c.baseUrl, strings.Replace(name, " ", "+", -1))

	c.log.InfoContext(ctx, "Searching City", "cityName", name)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var cityList CityList
	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&cityList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}
	if len(cityList.Cities) == 0 {
		return nil, fmt.Errorf("City not found for name: %s", name)
	}

	return &cityList.Cities[0], nil
}

func (c *ForecastApiClient) GetCityForecast(ctx context.Context, cityID string) (*Forecast, error) {
	url := fmt.Sprintf("%s/cidade/%s/previsao.xml", c.baseUrl, cityID)
	c.log.InfoContext(ctx, "Getting weather forecast", "cityId", cityID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var forecast Forecast
	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&forecast)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &forecast, nil
}

func (c *ForecastApiClient) GetWaveForecast(ctx context.Context, cityID string, day string) (*WaveForecast, error) {
	url := fmt.Sprintf("%s/cidade/%s/dia/%s/ondas.xml", c.baseUrl, cityID, day)

	c.log.InfoContext(ctx, "Getting wave forecast", "cityId", cityID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wave forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var waveForecast WaveForecast
	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&waveForecast)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &waveForecast, nil
}
