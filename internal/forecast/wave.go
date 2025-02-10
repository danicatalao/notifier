package forecast

import "encoding/xml"

type Wave struct {
	XMLName  xml.Name `xml:"cidade"`
	CityCode string   `xml:"codigo"`
	Name     string   `xml:"nome"`
	State    string   `xml:"uf"`
	Updated  string   `xml:"atualizacao"`
	Waves    []struct {
		Date      string `xml:"dia"`
		Wave      string `xml:"altura"`
		Direction string `xml:"direcao"`
		Agitation string `xml:"agitacao"`
	} `xml:"previsao"`
}
