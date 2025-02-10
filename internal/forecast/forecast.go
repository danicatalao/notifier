package forecast

import "encoding/xml"

type Forecast struct {
	XMLName   xml.Name `xml:"cidade"`
	Name      string   `xml:"nome"`
	State     string   `xml:"uf"`
	Updated   string   `xml:"atualizacao"`
	Forecasts []struct {
		Date  string `xml:"dia"`
		Time  string `xml:"hora"`
		Temp  string `xml:"temperatura"`
		Rain  string `xml:"iuv"`
		Cloud string `xml:"nebulosidade"`
	} `xml:"previsao"`
}
