package service

import (
	"context"
	general_dto "tcg_card_battler/web-api/internal/dto/general"
	"tcg_card_battler/web-api/internal/repository"
)

type GeneralService interface {
	GetElements(ctx context.Context) ([]general_dto.Element, error)
	GetOrigins(ctx context.Context) ([]general_dto.Origin, error)
}

type generalServiceImpl struct {
	generalRepo repository.GeneralRepository
}

func NewGeneralService(g repository.GeneralRepository) GeneralService {
	return &generalServiceImpl{generalRepo: g}
}

func (s *generalServiceImpl) GetElements(ctx context.Context) ([]general_dto.Element, error) {
	return s.generalRepo.GetElements(ctx)
}

func (s *generalServiceImpl) GetOrigins(ctx context.Context) ([]general_dto.Origin, error) {
	return s.generalRepo.GetOrigins(ctx)
}
