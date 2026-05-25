package store_dto

import booster_dto "tcg_card_battler/web-api/internal/dto/booster"

type PostBuyBoosterPackRQ struct {
	BoosterCode string `json:"boosterCode"`
	QTY         int    `json:"qty"`
}

type PostBuyBoosterPackRS struct {
	Cards [][]booster_dto.BoosterCard `json:"cards"`
}
