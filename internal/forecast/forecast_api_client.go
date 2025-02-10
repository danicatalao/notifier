package forecast

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type HttpClient interface {
	Get(s string) (resp *http.Response, err error)
}

type ForecastApiClient struct {
	httpClient HttpClient
	baseUrl    string
}

func NewForecastApiClient(c HttpClient, s string) *ForecastApiClient {
	return &ForecastApiClient{
		httpClient: c,
		baseUrl:    s,
	}
}

func (c *ForecastApiClient) GetCityForecast(cityID string) (*Forecast, error) {
	url := fmt.Sprintf("%s/cidade/%s/previsao.xml", c.baseUrl, cityID)

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
	if err := xml.Unmarshal(body, &forecast); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &forecast, nil
}
