package service

import (
	"context"
	account_dto "tcg_card_battler/web-api/internal/dto/account"
	"tcg_card_battler/web-api/internal/model"
	"tcg_card_battler/web-api/internal/repository"
)

type AccountService interface {
	GetAccountByEmail(ctx context.Context, email string) (*model.Account, error)
	GetAccountByID(ctx context.Context, accountID string) (*account_dto.AccountDetailRS, error)
	UpdateAccountGold(ctx context.Context, accountID string, gold int64) error
}

type accountServiceImpl struct {
	AccountRepo repository.AccountRepository
}

func NewAccountService(r repository.AccountRepository) AccountService {
	return &accountServiceImpl{AccountRepo: r}
}

func (s *accountServiceImpl) GetAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	return s.AccountRepo.GetAccountByEmail(ctx, email)
}

func (s *accountServiceImpl) GetAccountByID(ctx context.Context, accountID string) (*account_dto.AccountDetailRS, error) {
	return s.AccountRepo.GetAccountByID(ctx, accountID)
}

func (s *accountServiceImpl) UpdateAccountGold(ctx context.Context, accountID string, gold int64) error {
	_, err := s.AccountRepo.UpdateGold(ctx, accountID, gold)
	return err
}
