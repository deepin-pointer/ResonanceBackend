package model

type Goods struct {
	Name         string `json:"name"`
	Note         string `json:"note"`
	IsSpecial    bool   `json:"special"`
	OriginCityID int    `json:"origin"`
	BaseCount    int    `json:"count"`
	BasePrice    []int  `json:"base"`
}

type City struct {
	Name     string `json:"name"`
	Note     string `json:"note"`
	Distance []int  `json:"distance"`
}

type PriceRecord struct {
	GoodsID int `json:"goods"`
	CityID  int `json:"city"`
	Price   int `json:"price"`
}
