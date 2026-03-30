package booster_dto

type GetAllBoosterRS struct {
	Boosters []Booster `json:"boosters"`
}

type Booster struct {
	BoosterCode string `json:"boosterCode"`
	BoosterName string `json:"boosterName"`
	Price       int    `json:"price"`
}

type GetAllBoosterCardRS struct {
	Cards []BoosterCard `json:"cards"`
}

type BoosterCard struct {
	CardCode        string `json:"cardCode"`
	ImageTypeNumber int    `json:"imageTypeNumber"`
	CardTypeCode    string `json:"cardTypeCode"`
	CardRarityCode  string `json:"cardRarityCode"`
	Price           int    `json:"price"`
	Origin          string `json:"origin"`
}

type GetBoosterRarityRateRS struct {
	Items []BoosterCardPercentage `json:"items"`
}

type BoosterCardPercentage struct {
	CardRarityCode string `json:"cardRarityCode"`
	Percentage     int    `json:"percentage"`
}
