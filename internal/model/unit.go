package model

type Unit struct {
	UnitCode  string `db:"unit_code"`
	UnitName  string `db:"unit_name"`
	Origin    string `db:"origin"`
	UnitLevel int    `db:"unit_level"`
	Offense   int    `json:"offense"`
	Defense   int    `json:"defense"`
	Technique int    `json:"technique"`
	Speed     int    `json:"speed"`
	Spirit    int    `json:"spirit"`
}
