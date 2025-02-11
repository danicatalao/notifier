package forecast

type ForecastWave struct {
	Forecast *Forecast     `json:"previsão_do_tempo"`
	Wave     *WaveForecast `json:"ondas_do_dia,omitempty"`
}
