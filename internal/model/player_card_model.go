package model

type PlayerCard struct {
	AccountID       string `db:"account_id"`
	CardCode        string `db:"card_code"`
	ImageTypeNumber int    `db:"image_type_number"`
	QTY             int    `db:"qty"`
}
