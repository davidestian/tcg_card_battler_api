package service

import (
	"context"
	"errors"
	"fmt"
	inv_dto "tcg_card_battler/web-api/internal/dto/inventory"
	"tcg_card_battler/web-api/internal/repository"
)

type InventoryService interface {
	GetPlayerUnits(ctx context.Context, accountID string, limit, page int) (*inv_dto.GetPlayerUnitRS, error)
	GetPlayerUnitDetailByCode(ctx context.Context, accountID string, playerUnitID string) (*inv_dto.PlayerUnitDetailRS, error)
	GetAllPlayerCards(ctx context.Context, accountID string, limit int, cursorPrice int, cursorCardCode string, cursorImageTypeNumber int, pageNumber int, isPrev bool) (*inv_dto.GetAllPlayerCardRS, error)
	GetPlayerUnitCardByUnitCode(ctx context.Context, accountID string, unitCode string) ([]inv_dto.PlayerCard, error)
	PostPlayerUnitLevelUp(ctx context.Context, accountID string, rq inv_dto.PostPlayerUnitLevelUpRQ) error
	GetPlayerUnitPrevLevel(ctx context.Context, accountID string, playerUnitID string) ([]inv_dto.PlayerUnitPrevLevelRS, error)
	ChangePlayerUnitImage(ctx context.Context, accountID string, rq inv_dto.PlayerUnitLevelChangeImageRQ) error
	PostPlayerUnitUpgrade(ctx context.Context, accountID string, rq inv_dto.PostPlayerUnitUpgradeRQ) error
	GetEligibleUnitsToCreate(ctx context.Context, accountID string, limit, page int) (*inv_dto.GetEligibleUnitsToCreateRS, error)
	PostCreatePlayerUnit(ctx context.Context, accountID string, rq inv_dto.PostCreatePlayerUnitRQ) error
}

type inventoryServiceImpl struct {
	inventoryRepo repository.InventoryRepository
	unitRepo      repository.UnitRepository
	transactor    repository.Transactor
}

func NewInventoryService(inv repository.InventoryRepository, u repository.UnitRepository, trans repository.Transactor) InventoryService {
	return &inventoryServiceImpl{inventoryRepo: inv, unitRepo: u, transactor: trans}
}

func (h *inventoryServiceImpl) GetPlayerUnits(ctx context.Context, accountID string, limit, page int) (*inv_dto.GetPlayerUnitRS, error) {
	offset := (page - 1) * limit
	totalRecords, err := h.inventoryRepo.GetPlayerUnitCount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	units, err := h.inventoryRepo.GetPlayerUnits(ctx, accountID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalPage := 0
	if limit > 0 {
		totalPage = (totalRecords + limit - 1) / limit
	}

	return &inv_dto.GetPlayerUnitRS{
		TotalPage: totalPage,
		Units:     units,
	}, nil
}

func (h *inventoryServiceImpl) GetPlayerUnitDetailByCode(ctx context.Context, accountID string, playerUnitID string) (*inv_dto.PlayerUnitDetailRS, error) {
	return h.inventoryRepo.InvGetPlayerUnitDetailByID(ctx, accountID, playerUnitID)
}

func (h *inventoryServiceImpl) GetAllPlayerCards(ctx context.Context, accountID string, limit int, cursorPrice int, cursorCardCode string, cursorImageTypeNumber int, pageNumber int, isPrev bool) (*inv_dto.GetAllPlayerCardRS, error) {
	return h.inventoryRepo.GetAllPlayerCards(ctx, accountID, limit, cursorPrice, cursorCardCode, cursorImageTypeNumber, pageNumber, isPrev)
}

func (h *inventoryServiceImpl) GetPlayerUnitCardByUnitCode(ctx context.Context, accountID string, unitCode string) ([]inv_dto.PlayerCard, error) {
	return h.inventoryRepo.InvGetPlayerUnitCardByUnitCode(ctx, accountID, unitCode)
}

func (s *inventoryServiceImpl) PostPlayerUnitLevelUp(ctx context.Context, accountID string, rq inv_dto.PostPlayerUnitLevelUpRQ) error {
	if len(rq.Items) == 0 {
		return fmt.Errorf("no items selected")
	}

	playerUnit, err := s.inventoryRepo.InvGetPlayerUnitDetailByID(ctx, accountID, rq.PlayerUnitID)
	if err != nil {
		return err
	}
	if playerUnit.PlayerUnitLevel >= 10 {
		return fmt.Errorf("unit already max level")
	}

	totalRequest := 0
	itns := make([]int, len(rq.Items))
	for i, item := range rq.Items {
		totalRequest += item.QTY
		itns[i] = item.ImageTypeNumber
	}

	if totalRequest < 10 {
		return fmt.Errorf("not enough resources: need %d, got %d", playerUnit.PlayerUnitLevel, totalRequest)
	}

	return s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		cards, err := s.inventoryRepo.GetPlayerCardByCodeAndTypeNumber(txCtx, accountID, playerUnit.FirstUnitCode, itns)
		if err != nil {
			return err
		}

		updateIDs := make([]int, 0, len(rq.Items))
		updateQTYs := make([]int, 0, len(rq.Items))
		deleteIDs := make([]int, 0, len(rq.Items))
		for _, item := range rq.Items {
			val, exist := cards[item.ImageTypeNumber]
			if !exist || val < item.QTY {
				return fmt.Errorf("insufficient cards for type %d", item.ImageTypeNumber)
			}

			finalQTY := val - item.QTY
			if finalQTY == 0 {
				deleteIDs = append(deleteIDs, item.ImageTypeNumber)
			} else {
				updateIDs = append(updateIDs, item.ImageTypeNumber)
				updateQTYs = append(updateQTYs, finalQTY)
			}
		}

		if err := s.inventoryRepo.BatchUpdatePlayerCards(txCtx, accountID, playerUnit.FirstUnitCode, updateIDs, updateQTYs); err != nil {
			return err
		}

		if err := s.inventoryRepo.BatchDeletePlayerCard(txCtx, accountID, playerUnit.FirstUnitCode, deleteIDs); err != nil {
			return err
		}

		if err := s.inventoryRepo.IncrementUnitLevel(txCtx, rq.PlayerUnitID); err != nil {
			return err
		}

		return nil
	})
}

func (h *inventoryServiceImpl) GetPlayerUnitPrevLevel(ctx context.Context, accountID string, playerUnitID string) ([]inv_dto.PlayerUnitPrevLevelRS, error) {
	return h.inventoryRepo.InvGetPlayerUnitPrevLevel(ctx, accountID, playerUnitID)
}

func (s *inventoryServiceImpl) ChangePlayerUnitImage(ctx context.Context, accountID string, rq inv_dto.PlayerUnitLevelChangeImageRQ) error {
	return s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		newQty, err := s.inventoryRepo.DecrementCard(txCtx, accountID, rq.UnitCode, rq.ImageTypeNumber)
		if err != nil {
			return err
		}

		if newQty == 0 {
			if err := s.inventoryRepo.DeleteCard(txCtx, accountID, rq); err != nil {
				return err
			}
		}

		return s.inventoryRepo.UpdateUnitLevelImage(txCtx, accountID, rq)
	})
}

func (s *inventoryServiceImpl) PostPlayerUnitUpgrade(ctx context.Context, accountID string, rq inv_dto.PostPlayerUnitUpgradeRQ) error {
	playerUnit, err := s.inventoryRepo.InvGetPlayerUnitDetailByID(ctx, accountID, rq.PlayerUnitID)
	if err != nil {
		return errors.New("player unit not found")
	}

	unitLevelPath, err := s.unitRepo.GetUnitLevelPathByCode(ctx, playerUnit.LastUnitCode, rq.TargetUnitCode)
	if err != nil {
		return errors.New("unit next level not found")
	}

	totalRequest := 0
	itns := make([]int, len(rq.Items))
	for i, item := range rq.Items {
		totalRequest += item.QTY
		itns[i] = item.ImageTypeNumber
	}

	if totalRequest < 20 {
		return errors.New("not enough resources for upgrade")
	}

	return s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		cards, err := s.inventoryRepo.GetPlayerCardByCodeAndTypeNumber(ctx, accountID, unitLevelPath.UnitCode, itns)
		if err != nil {
			return err
		}

		updateIDs := make([]int, 0, len(rq.Items))
		updateQTYs := make([]int, 0, len(rq.Items))
		deleteIDs := make([]int, 0, len(rq.Items))
		for _, item := range rq.Items {
			val, exist := cards[item.ImageTypeNumber]
			if !exist || val < item.QTY {
				return fmt.Errorf("insufficient cards for type %d", item.ImageTypeNumber)
			}

			finalQTY := val - item.QTY
			if finalQTY == 0 {
				deleteIDs = append(deleteIDs, item.ImageTypeNumber)
			} else {
				updateIDs = append(updateIDs, item.ImageTypeNumber)
				updateQTYs = append(updateQTYs, finalQTY)
			}
		}

		if err := s.inventoryRepo.BatchUpdatePlayerCards(ctx, accountID, playerUnit.FirstUnitCode, updateIDs, updateQTYs); err != nil {
			return err
		}

		if err := s.inventoryRepo.BatchDeletePlayerCard(ctx, accountID, playerUnit.FirstUnitCode, deleteIDs); err != nil {
			return err
		}

		if err := s.inventoryRepo.InsertPlayerLevel(ctx, rq.PlayerUnitID, unitLevelPath.UnitLevel, unitLevelPath.UnitCode); err != nil {
			return err
		}
		return nil
	})
}

func (s *inventoryServiceImpl) GetEligibleUnitsToCreate(ctx context.Context, accountID string, limit, page int) (*inv_dto.GetEligibleUnitsToCreateRS, error) {
	offset := (page - 1) * limit

	totalRecords, err := s.inventoryRepo.GetEligibleUnitsCount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	units, err := s.inventoryRepo.GetEligibleUnitsList(ctx, accountID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalPage := 0
	if limit > 0 {
		totalPage = (totalRecords + limit - 1) / limit
	}

	return &inv_dto.GetEligibleUnitsToCreateRS{
		TotalPage: totalPage,
		CurrPage:  page,
		Units:     units,
	}, nil
}

func (s *inventoryServiceImpl) PostCreatePlayerUnit(ctx context.Context, accountID string, rq inv_dto.PostCreatePlayerUnitRQ) error {
	if len(rq.Items) == 0 {
		return fmt.Errorf("no items selected")
	}

	totalRequest := 0
	itns := make([]int, len(rq.Items))
	for i, item := range rq.Items {
		totalRequest += item.QTY
		itns[i] = item.ImageTypeNumber
	}

	if totalRequest < 50 {
		return fmt.Errorf("insufficient selection: total QTY must be at least 50, got %d", totalRequest)
	}

	return s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		cards, err := s.inventoryRepo.GetPlayerCardByCodeAndTypeNumber(txCtx, accountID, rq.UnitCode, itns)
		if err != nil {
			return err
		}

		updateIDs := make([]int, 0, len(rq.Items))
		updateQTYs := make([]int, 0, len(rq.Items))
		deleteIDs := make([]int, 0, len(rq.Items))

		for _, item := range rq.Items {
			currentQTY, exist := cards[item.ImageTypeNumber]
			if !exist || currentQTY < item.QTY {
				return fmt.Errorf("insufficient cards for type %d (owned: %d, req: %d)",
					item.ImageTypeNumber, currentQTY, item.QTY)
			}

			finalQTY := currentQTY - item.QTY
			if finalQTY == 0 {
				deleteIDs = append(deleteIDs, item.ImageTypeNumber)
			} else {
				updateIDs = append(updateIDs, item.ImageTypeNumber)
				updateQTYs = append(updateQTYs, finalQTY)
			}
		}

		if len(updateIDs) > 0 {
			if err := s.inventoryRepo.BatchUpdatePlayerCards(txCtx, accountID, rq.UnitCode, updateIDs, updateQTYs); err != nil {
				return fmt.Errorf("failed to update cards: %w", err)
			}
		}

		if len(deleteIDs) > 0 {
			if err := s.inventoryRepo.BatchDeletePlayerCard(txCtx, accountID, rq.UnitCode, deleteIDs); err != nil {
				return fmt.Errorf("failed to delete empty cards: %w", err)
			}
		}

		playerUnitID, err := s.inventoryRepo.InsertPlayerUnit(txCtx, accountID)
		if err != nil {
			return fmt.Errorf("failed to create unit: %w", err)
		}

		if err = s.inventoryRepo.InsertPlayerLevel(txCtx, playerUnitID, 1, rq.UnitCode); err != nil {
			return fmt.Errorf("failed to set unit level: %w", err)
		}

		return nil
	})
}
