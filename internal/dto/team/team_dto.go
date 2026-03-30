package team_dto

import (
	inv_dto "tcg_card_battler/web-api/internal/dto/inventory"

	"github.com/gofrs/uuid/v5"
)

type GetPlayerTeamRS struct {
	PlayerTeams []PlayerTeam `json:"playerTeams"`
	TotalPage   int          `json:"totalPage"`
}

type PlayerTeam struct {
	TeamID      uuid.UUID          `json:"teamID"`
	TeamName    string             `json:"teamName"`
	IsActive    bool               `json:"isActive"`
	TeamLevel   int                `json:"teamLevel"`
	PlayerUnit1 inv_dto.PlayerUnit `json:"playerUnit1"`
	PlayerUnit2 inv_dto.PlayerUnit `json:"playerUnit2"`
	PlayerUnit3 inv_dto.PlayerUnit `json:"playerUnit3"`
}

type PostPlayerTeamRQ struct {
	TeamName      string `json:"teamName"`
	PlayerUnitID1 string `json:"playerUnitID1"`
	PlayerUnitID2 string `json:"playerUnitID2"`
	PlayerUnitID3 string `json:"playerUnitID3"`
}

type PutActivePlayerTeamRQ struct {
	TeamID string `json:"teamID"`
}

type PlayerTeamUnit struct {
	PlayerUnitID    string   `json:"playerUnitID"`
	PlayerUnitLevel int      `json:"playerUnitLevel"`
	ImageTypeNumber int      `json:"imageTypeNumber"`
	UnitCode        string   `json:"unitCode"`
	UnitName        string   `json:"unitName"`
	Tags            []string `json:"tags"`
	Origin          string   `json:"origin"`
	UnitLevel       int      `json:"unitLevel"`
	Offense         int      `json:"offense"`
	Defense         int      `json:"defense"`
	Technique       int      `json:"technique"`
	Speed           int      `json:"speed"`
	Spirit          int      `json:"spirit"`
	ElementID1      int      `json:"elementID1"`
	ElementID2      int      `json:"elementID2"`
}
