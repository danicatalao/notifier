package forecast

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html/charset"
)

type HttpClient interface {
	Get(s string) (resp *http.Response, err error)
}

type ForecastApiClient struct {
	httpClient HttpClient
	baseUrl    string
}

func NewForecastApiClient(c HttpClient, s string) ForecastApiClient {
	return ForecastApiClient{
		httpClient: c,
		baseUrl:    s,
	}
}

func (c *ForecastApiClient) GetCity(name string) (*City, error) {
	url := fmt.Sprintf("%s/listaCidades?city=%s", c.baseUrl, name)

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

	return &cityList.Cities[0], nil
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
	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&forecast)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &forecast, nil
}

func (c *ForecastApiClient) GetWaveForecast(cityID string, day string) (*WaveForecast, error) {
	url := fmt.Sprintf("%s/cidade/%s/dia/%s/ondas.xml", c.baseUrl, cityID, day)

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
