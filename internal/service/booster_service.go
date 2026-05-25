package service

import (
	"context"
	booster_dto "tcg_card_battler/web-api/internal/dto/booster"
	"tcg_card_battler/web-api/internal/repository"
)

type BoosterService interface {
	GetAllBooster(ctx context.Context, limit int) (*booster_dto.GetAllBoosterRS, error)
	GetAllBoosterCard(ctx context.Context, boosterCode string) (*booster_dto.GetAllBoosterCardRS, error)
	GetBoosterRarityRate(ctx context.Context, boosterCode string) (*booster_dto.GetBoosterRarityRateRS, error)
}

type boosterServiceImpl struct {
	boosterRepo repository.BoosterRepository
}

func NewBoosterService(r repository.BoosterRepository) BoosterService {
	return &boosterServiceImpl{boosterRepo: r}
}

func (s *boosterServiceImpl) GetAllBooster(ctx context.Context, limit int) (*booster_dto.GetAllBoosterRS, error) {
	data, err := s.boosterRepo.GetAllBooster(ctx)
	if err != nil {
		return nil, err
	}

	for idx, booster := range data.Boosters {
		cards, err := s.boosterRepo.GetHighestBoosterCard(ctx, booster.BoosterCode, 5)

		if err != nil {
			return nil, err
		}
		data.Boosters[idx].BoosterCards = cards.Cards
	}

	return data, nil
}

func (s *boosterServiceImpl) GetAllBoosterCard(ctx context.Context, boosterCode string) (*booster_dto.GetAllBoosterCardRS, error) {
	return s.boosterRepo.GetAllBoosterCard(ctx, boosterCode)
}

func (s *boosterServiceImpl) GetBoosterRarityRate(ctx context.Context, boosterCode string) (*booster_dto.GetBoosterRarityRateRS, error) {
	return s.boosterRepo.GetBoosterRarityRate(ctx, boosterCode)
}
