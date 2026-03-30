package account_dto

type AccountDetailRS struct {
	AccountID   string `json:"accountID"`
	AccountName string `json:"accountName"`
	Email       string `json:"email"`
	Gold        int64  `json:"gold"`
}

type PutAccountGoldRQ struct {
	Gold int64 `json:"gold"`
}
