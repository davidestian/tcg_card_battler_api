package booster_dto

type GetAllBoosterRS struct {
	Boosters []Booster `json:"boosters"`
}

type Booster struct {
	BoosterCode  string        `json:"boosterCode"`
	BoosterName  string        `json:"boosterName"`
	Price        int           `json:"price"`
	BoosterCards []BoosterCard `json:"boosterCards"`
}

type GetAllBoosterCardRS struct {
	Cards []BoosterCard `json:"cards"`
}

type BoosterCard struct {
	CardCode        string `json:"cardCode"`
	ImageTypeNumber int    `json:"imageTypeNumber"`
	CardRarityCode  string `json:"cardRarityCode"`
	Price           int    `json:"price"`
	Origin          string `json:"origin"`
	ElementID1      int    `json:"elementID1"`
	ElementID2      int    `json:"elementID2"`
}

type GetBoosterRarityRateRS struct {
	Items []BoosterCardPercentage `json:"items"`
}

type BoosterCardPercentage struct {
	CardRarityCode string `json:"cardRarityCode"`
	Percentage     int    `json:"percentage"`
}
