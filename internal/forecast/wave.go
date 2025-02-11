package forecast

type WaveForecast struct {
	Name          string `xml:"nome" json:"nome"`
	State         string `xml:"uf" json:"uf"`
	Updated       string `xml:"atualizacao" json:"atualizacao"`
	MorningWave   Wave   `xml:"manha" json:"manha"`
	AfternoonWave Wave   `xml:"tarde" json:"tarde"`
	NightWave     Wave   `xml:"noite" json:"noite"`
}

type Wave struct {
	Date          string `xml:"dia" json:"dia"`
	Agitation     string `xml:"agitacao" json:"agitacao"`
	Height        string `xml:"altura" json:"altura"`
	Direction     string `xml:"direcao" json:"direcao"`
	Wind          string `xml:"vento" json:"vento"`
	WindDirection string `xml:"vento_dir" json:"vento_dir"`
}
