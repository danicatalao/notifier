package forecast

type CityList struct {
	Cities []City `xml:"cidade"`
}

type City struct {
	Id    int    `xml:"id"`
	Name  string `xml:"nome"`
	State string `xml:"uf"`
}
