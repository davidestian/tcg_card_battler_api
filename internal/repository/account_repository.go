package repository

import (
	"context"
	"errors"
	account_dto "tcg_card_battler/web-api/internal/dto/account"
	"tcg_card_battler/web-api/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	exec(ctx context.Context) DBQuerier
	Insert(ctx context.Context, email, username, passwordHash string) (string, error)
	GetAccountByEmail(ctx context.Context, email string) (*model.Account, error)
	GetAccountByID(ctx context.Context, accountID string) (*account_dto.AccountDetailRS, error)
	UpdateGold(ctx context.Context, accountID string, gold int64) (int, error)
	ForgotPassword(ctx context.Context, email, passwordHash string) error
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

func (r *AccountRepositoryImpl) Insert(ctx context.Context, email, username, passwordHash string) (string, error) {
	query := `
    INSERT INTO public.accounts (account_id, email, account_name, password_hash, gold)
    VALUES (gen_random_uuid(), $1, $2, $3, 0)
    ON CONFLICT (email) DO NOTHING
	RETURNING account_id;
	`
	var accountID string
	err := r.exec(ctx).QueryRow(ctx, query, email, username, passwordHash).Scan(&accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("email already exist")
		}
		return "", err
	}

	return accountID, nil
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

func (r *AccountRepositoryImpl) ForgotPassword(ctx context.Context, email, passwordHash string) error {
	query := `UPDATE accounts SET 
			password_hash = $2 
		WHERE email = $1;
	`
	_, err := r.exec(ctx).Exec(ctx, query, email, passwordHash)
	return err
}
