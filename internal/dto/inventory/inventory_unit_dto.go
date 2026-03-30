package inv_dto

type PlayerUnitDetailRS struct {
	PlayerUnitID    string   `json:"playerUnitID"`
	FirstUnitCode   string   `json:"firstUnitCode"`
	LastUnitCode    string   `json:"lastUnitCode"`
	Origin          string   `json:"origin"`
	UnitName        string   `json:"unitName"`
	PlayerUnitLevel int      `json:"playerUnitLevel"`
	Tags            []string `json:"tags"`
	ImageTypeNumber int      `json:"imageTypeNumber"`
}

type GetPlayerUnitRS struct {
	TotalPage int          `json:"totalPage"`
	Units     []PlayerUnit `json:"units"`
}

type PlayerUnit struct {
	PlayerUnitID    string `json:"playerUnitID"`
	Level           int    `json:"level"`
	UnitCode        string `json:"unitCode"`
	Origin          string `json:"origin"`
	ImageTypeNumber int    `json:"imageTypeNumber"`
	ElementID1      int    `json:"elementID1"`
	ElementID2      int    `json:"elementID2"`
}

type playerUnitLevelUpRS struct {
	PlayerUnitID    string `json:"playerUnitID"`
	Level           int    `json:"level"`
	UnitCode        string `json:"unitCode"`
	Origin          string `json:"origin"`
	ImageTypeNumber int    `json:"imageTypeNumber"`
}

type PlayerUnitPrevLevelRS struct {
	PlayerUnitID    string `json:"playerUnitID"`
	TargetLevel     int    `json:"targetLevel"`
	ImageTypeCount  int    `json:"imageTypeCount"`
	UnitName        string `json:"unitName"`
	UnitCode        string `json:"unitCode"`
	Origin          string `json:"origin"`
	ImageTypeNumber int    `json:"imageTypeNumber"`
	Offense         int    `json:"offense"`
	Defense         int    `json:"defense"`
	Technique       int    `json:"technique"`
	Speed           int    `json:"speed"`
	Spirit          int    `json:"spirit"`
}

type PlayerUnitLevelChangeImageRQ struct {
	PlayerUnitID    string `json:"playerUnitID"`
	TargetLevel     int    `json:"targetLevel"`
	UnitCode        string `json:"unitCode"`
	ImageTypeNumber int    `json:"imageTypeNumber"`
}

type PostPlayerUnitLevelUpRQ struct {
	PlayerUnitID string           `json:"playerUnitID"`
	Items        []CostCardItemRQ `json:"items"`
}

type PostCreatePlayerUnitRQ struct {
	UnitCode string           `json:"unitCode"`
	Items    []CostCardItemRQ `json:"items"`
}

type PostPlayerUnitUpgradeRQ struct {
	PlayerUnitID   string           `json:"playerUnitID"`
	TargetUnitCode string           `json:"targetUnitCode"`
	Items          []CostCardItemRQ `json:"items"`
}

type CostCardItemRQ struct {
	ImageTypeNumber int `json:"imageTypeNumber"`
	QTY             int `json:"qty"`
}

type GetAllPlayerCardRS struct {
	TotalPage int          `json:"totalPage"`
	CurrPage  int          `json:"currPage"`
	Cards     []PlayerCard `json:"cards"`
}

type PlayerCard struct {
	CardCode        string `json:"cardCode"`
	ImageTypeNumber int    `json:"imageTypeNumber"`
	Origin          string `json:"origin"`
	QTY             int    `json:"qty"`
	CardRarityCode  string `json:"cardRarityCode"`
	Price           int    `json:"price"`
}

type GetEligibleUnitsToCreateRS struct {
	TotalPage int            `json:"totalPage"`
	CurrPage  int            `json:"currPage"`
	Units     []EligibleUnit `json:"units"`
}
type EligibleUnit struct {
	UnitCode string `json:"unitCode"`
	Origin   string `json:"origin"`
}
