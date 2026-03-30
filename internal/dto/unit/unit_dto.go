package unit_dto

type GetUnitNextLevelPathRS struct {
	UnitCode   string   `json:"unitCode"`
	UnitName   string   `json:"unitName"`
	Tags       []string `json:"tags"`
	Origin     string   `json:"origin"`
	UnitLevel  int      `json:"unitLevel"`
	Offense    int      `json:"offense"`
	Defense    int      `json:"defense"`
	Technique  int      `json:"technique"`
	Speed      int      `json:"speed"`
	Spirit     int      `json:"spirit"`
	ElementID1 int      `json:"elementID1"`
	ElementID2 int      `json:"elementID2"`
}

type Unit struct {
	UnitCode       string   `json:"unitCode"`
	UnitName       string   `json:"unitName"`
	Tags           []string `json:"tags"`
	Origin         string   `json:"origin"`
	UnitLevel      int      `json:"unitLevel"`
	ImageTypeCount int      `json:"imageTypeCount"`
	Offense        int      `json:"offense"`
	Defense        int      `json:"defense"`
	Technique      int      `json:"technique"`
	Speed          int      `json:"speed"`
	Spirit         int      `json:"spirit"`
	ElementID1     int      `json:"elementID1"`
	ElementID2     int      `json:"elementID2"`
}
