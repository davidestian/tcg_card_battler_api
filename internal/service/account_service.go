package service

import (
	"context"
	account_dto "tcg_card_battler/web-api/internal/dto/account"
	"tcg_card_battler/web-api/internal/model"
	"tcg_card_battler/web-api/internal/repository"

	"github.com/alexedwards/argon2id"
)

type AccountService interface {
	GetAccountByEmail(ctx context.Context, email string) (*model.Account, error)
	GetAccountByID(ctx context.Context, accountID string) (*account_dto.AccountDetailRS, error)
	UpdateAccountGold(ctx context.Context, accountID string, gold int64) error
	CreateNewAccount(ctx context.Context, email, username, password string) error
	ForgotPassword(ctx context.Context, email, password string) error
}

type accountServiceImpl struct {
	accountRepo   repository.AccountRepository
	unitRepo      repository.UnitRepository
	inventoryRepo repository.InventoryRepository
	teamRepo      repository.TeamRepository
	transactor    repository.Transactor
}

func NewAccountService(r repository.AccountRepository, u repository.UnitRepository, i repository.InventoryRepository, tr repository.TeamRepository, t repository.Transactor) AccountService {
	return &accountServiceImpl{accountRepo: r, unitRepo: u, inventoryRepo: i, teamRepo: tr, transactor: t}
}

func (s *accountServiceImpl) GetAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	return s.accountRepo.GetAccountByEmail(ctx, email)
}

func (s *accountServiceImpl) GetAccountByID(ctx context.Context, accountID string) (*account_dto.AccountDetailRS, error) {
	return s.accountRepo.GetAccountByID(ctx, accountID)
}

func (s *accountServiceImpl) UpdateAccountGold(ctx context.Context, accountID string, gold int64) error {
	_, err := s.accountRepo.UpdateGold(ctx, accountID, gold)
	return err
}

func (s *accountServiceImpl) CreateNewAccount(ctx context.Context, email, username, password string) error {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		accountID, err := s.accountRepo.Insert(txCtx, email, username, hash)
		if err != nil {
			return err
		}

		level := 9
		playerUnitIDs := make([]string, 3)
		n := 0
		for n < 3 {
			units, err := s.unitRepo.GetRandomUnitByLevel(txCtx, level)
			if err != nil {
				return err
			}

			playerUnitID, err := s.inventoryRepo.InsertPlayerUnit(txCtx, accountID, level)
			if err != nil {
				return err
			}

			for _, unit := range units {
				if err = s.inventoryRepo.InsertPlayerLevel(txCtx, playerUnitID, unit.UnitLevel, unit.UnitCode); err != nil {
					return err
				}
			}
			playerUnitIDs[n] = playerUnitID
			n++
		}

		err = s.teamRepo.InsertPlayerTeam(txCtx, accountID, "TEAM 1", playerUnitIDs[0], playerUnitIDs[1], playerUnitIDs[2], "1")
		return err
	})

	return err
}

func (s *accountServiceImpl) ForgotPassword(ctx context.Context, email, password string) error {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	return s.accountRepo.ForgotPassword(ctx, email, hash)
}
