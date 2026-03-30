package model

import "github.com/gofrs/uuid/v5"

type PlayerUnit struct {
	PlayerUnitID    uuid.UUID `db:"player_unit_id"`
	AccountID       uuid.UUID `db:"account_id"`
	PlayerUnitLevel int       `db:"player_unit_level"`
}
