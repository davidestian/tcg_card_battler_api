package battle_dto

type BattleUnit struct {
	BattleUnitID string           `json:"battleUnitID"`
	Level        int              `json:"level"`
	Paths        []BattleUnitPath `json:"paths"`
}

type BattleUnitPath struct {
	UnitCode        string `json:"unitCode"`
	UnitName        string `json:"unitName"`
	UnitLevel       int    `json:"unitLevel"`
	Origin          string `json:"origin"`
	ImageTypeNumber int    `json:"imageTypeNumber"`
	Offense         int    `json:"offense"`
	Defense         int    `json:"defense"`
	Technique       int    `json:"technique"`
	Speed           int    `json:"speed"`
	Spirit          int    `json:"spirit"`
	ElementID1      int    `json:"elementID1"`
	ElementID2      int    `json:"elementID2"`
}
