package model

import "github.com/gofrs/uuid/v5"

type Account struct {
	AccountID    uuid.UUID `db:"account_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	AccountName  string    `db:"account_name"`
	Gold         int64     `db:"gold"`
}
