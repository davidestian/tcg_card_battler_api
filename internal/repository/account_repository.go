package repository

import (
	"context"
	account_dto "tcg_card_battler/web-api/internal/dto/account"
	"tcg_card_battler/web-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	exec(ctx context.Context) DBQuerier
	GetAccountByEmail(ctx context.Context, email string) (*model.Account, error)
	GetAccountByID(ctx context.Context, accountID string) (*account_dto.AccountDetailRS, error)
	UpdateGold(ctx context.Context, accountID string, gold int64) (int, error)
}

type AccountRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) AccountRepository {
	return &AccountRepositoryImpl{pool: pool}
}

func (r *AccountRepositoryImpl) exec(ctx context.Context) DBQuerier {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}

	return r.pool
}

func (r *AccountRepositoryImpl) GetAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	query := `SELECT account_id, email, account_name, password_hash FROM accounts WHERE email = $1`

	var accountData model.Account
	err := r.pool.QueryRow(ctx, query, email).Scan(&accountData.AccountID, &accountData.Email, &accountData.AccountName, &accountData.PasswordHash)
	return &accountData, err
}

func (r *AccountRepositoryImpl) GetAccountByID(ctx context.Context, accountID string) (*account_dto.AccountDetailRS, error) {
	query := `SELECT account_id, email, account_name, gold FROM accounts WHERE account_id = $1`

	var accountData account_dto.AccountDetailRS
	err := r.pool.QueryRow(ctx, query, accountID).Scan(&accountData.AccountID, &accountData.Email, &accountData.AccountName, &accountData.Gold)
	return &accountData, err
}

func (r *AccountRepositoryImpl) UpdateGold(ctx context.Context, accountID string, gold int64) (int, error) {
	var newGold int
	query := `UPDATE accounts SET 
			gold = gold + $2 
		WHERE account_id = $1 
		RETURNING gold
	`
	err := r.exec(ctx).QueryRow(ctx, query, accountID, gold).Scan(&newGold)
	return newGold, err
}
