package model

import "github.com/gofrs/uuid/v5"

type PlayerUnitLevel struct {
	PlayerUnitID  uuid.UUID `db:"player_unit_id"`
	TargetLevel   int       `db:"target_level"`
	UnitCode      string    `db:"unit_code"`
	ImageTypeCode string    `db:"image_type_code"`
}
