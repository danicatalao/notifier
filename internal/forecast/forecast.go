package forecast

type Forecast struct {
	Name      string `xml:"nome" json:"nome"`
	State     string `xml:"uf" json:"uf"`
	Updated   string `xml:"atualizacao" json:"atualizacao"`
	Forecasts []struct {
		Day       string  `xml:"dia" json:"dia"`
		Condition string  `xml:"tempo" json:"tempo"`
		Max       int     `xml:"maxima" json:"maxima"`
		Min       int     `xml:"minima" json:"minima"`
		Iuv       float32 `xml:"iuv" json:"iuv"`
	} `xml:"previsao" json:"previsao"`
}
