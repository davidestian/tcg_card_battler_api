package service

import (
	"context"
	"fmt"
	"math/rand/v2"
	booster_dto "tcg_card_battler/web-api/internal/dto/booster"
	store_dto "tcg_card_battler/web-api/internal/dto/store"
	"tcg_card_battler/web-api/internal/repository"
)

type StoreService interface {
	PostBuyBoosterPack(ctx context.Context, accountID string, rq store_dto.PostBuyBoosterPackRQ) (*store_dto.PostBuyBoosterPackRS, error)
}

type storeServiceImpl struct {
	accountRepo   repository.AccountRepository
	boosterRepo   repository.BoosterRepository
	inventoryRepo repository.InventoryRepository
	transactor    repository.Transactor
}

func NewStoreService(a repository.AccountRepository, r repository.BoosterRepository, inv repository.InventoryRepository, t repository.Transactor) StoreService {
	return &storeServiceImpl{accountRepo: a, boosterRepo: r, inventoryRepo: inv, transactor: t}
}

func (s *storeServiceImpl) PostBuyBoosterPack(ctx context.Context, accountID string, rq store_dto.PostBuyBoosterPackRQ) (*store_dto.PostBuyBoosterPackRS, error) {
	// 1. Fetch initial data
	account, err := s.accountRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	booster, err := s.boosterRepo.GetBoosterByCode(ctx, rq.BoosterCode)
	if err != nil {
		return nil, err
	}

	// 2. Validate Balance
	totalCost := int64(booster.Price * rq.QTY)
	if account.Gold < totalCost {
		return nil, fmt.Errorf("insufficient gold")
	}

	// 3. Prepare Gacha Data (Pre-grouping is key)
	rarities, err := s.boosterRepo.GetBoosterRarityRate(ctx, rq.BoosterCode)
	if err != nil {
		return nil, err
	}

	boosterCards, err := s.boosterRepo.GetAllBoosterCard(ctx, rq.BoosterCode)
	if err != nil {
		return nil, err
	}

	cardsByRarity := make(map[string][]booster_dto.BoosterCard)
	for _, card := range boosterCards.Cards {
		cardsByRarity[card.CardRarityCode] = append(cardsByRarity[card.CardRarityCode], card)
	}

	// Build cumulative percentage list
	cumulativeList := make([]int, len(rarities.Items))
	currentSum := 0
	for idx, item := range rarities.Items {
		currentSum += item.Percentage
		cumulativeList[idx] = currentSum
	}

	// 4. Generate Packs (Combined Loops)
	result := &store_dto.PostBuyBoosterPackRS{
		Cards: make([][]booster_dto.BoosterCard, rq.QTY),
	}

	playerCards := make(map[string]int)
	codes := make([]string, 0)
	imageTypes := make([]int32, 0)
	quantities := make([]int32, 0)
	tempIndex := 0

	for i := 0; i < rq.QTY; i++ {
		result.Cards[i] = make([]booster_dto.BoosterCard, 6)
		for j := 0; j < 6; j++ {
			roll := rand.IntN(100) + 1

			// Find Rarity
			var rarityCode string
			for idx, threshold := range cumulativeList {
				if roll <= threshold {
					rarityCode = rarities.Items[idx].CardRarityCode
					break
				}
			}

			// Select Random Card from Pool
			pool := cardsByRarity[rarityCode]
			if len(pool) == 0 {
				return nil, fmt.Errorf("no cards found for rarity: %s", rarityCode)
			}

			targetCard := pool[rand.IntN(len(pool))]
			result.Cards[i][j] = targetCard
			key := fmt.Sprintf("%s%d", targetCard.CardCode, targetCard.ImageTypeNumber)
			if val, exist := playerCards[key]; exist {
				quantities[val]++ // Increment the count for duplicates
			} else {
				codes = append(codes, targetCard.CardCode)
				imageTypes = append(imageTypes, int32(targetCard.ImageTypeNumber))
				quantities = append(quantities, 1)
				playerCards[key] = tempIndex
				tempIndex++
			}
		}
	}

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		resultGold, err := s.accountRepo.UpdateGold(txCtx, accountID, totalCost*-1)
		if err != nil {
			return err
		}
		if resultGold < 0 {
			return fmt.Errorf("insufficient gold")
		}

		err = s.inventoryRepo.BatchInsertPlayerCards(txCtx, accountID, codes, imageTypes, quantities)
		return err
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
