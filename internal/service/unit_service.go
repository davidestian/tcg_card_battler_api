package service

import (
	"context"
	unit_dto "tcg_card_battler/web-api/internal/dto/unit"
	"tcg_card_battler/web-api/internal/repository"
)

type UnitService interface {
	GetUnitByCode(ctx context.Context, unitCode string) (*unit_dto.Unit, error)
	GetUnitLevelPathByCode(ctx context.Context, unitCode string) ([]unit_dto.GetUnitNextLevelPathRS, error)
}
type unitServiceImpl struct {
	unitRepo repository.UnitRepository
}

func NewUnitService(r repository.UnitRepository) UnitService {
	return &unitServiceImpl{unitRepo: r}
}

func (s *unitServiceImpl) GetUnitByCode(ctx context.Context, unitCode string) (*unit_dto.Unit, error) {
	result, err := s.unitRepo.GetUnitByCode(ctx, unitCode)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *unitServiceImpl) GetUnitLevelPathByCode(ctx context.Context, unitCode string) ([]unit_dto.GetUnitNextLevelPathRS, error) {
	results, err := s.unitRepo.GetAllUnitLevelPathByCode(ctx, unitCode)
	if err != nil {
		return nil, err
	}

	return results, nil
}
