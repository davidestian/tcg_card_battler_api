package service

import (
	"context"
	"math/rand/v2"

	battle_dto "tcg_card_battler/web-api/internal/dto/battle"
	"tcg_card_battler/web-api/internal/repository"

	"github.com/gofrs/uuid/v5"
)

type BattleService interface {
	GetRandomEnemyBattleUnits(ctx context.Context, levels, evoLevels []int) ([]battle_dto.BattleUnit, error)
	GetPlayerTeamUnits(ctx context.Context, accountID, playerTeamID string) ([]battle_dto.BattleUnit, error)
}

type battleServiceImpl struct {
	unitRepo repository.UnitRepository
	teamRepo repository.TeamRepository
}

func NewBattleService(ur repository.UnitRepository, tr repository.TeamRepository) BattleService {
	return &battleServiceImpl{ur, tr}
}

func (s *battleServiceImpl) GetRandomEnemyBattleUnits(ctx context.Context, levels, evoLevels []int) ([]battle_dto.BattleUnit, error) {
	results := make([]battle_dto.BattleUnit, len(evoLevels))
	for index, val := range evoLevels {
		units, err := s.unitRepo.GetRandomUnitByLevel(ctx, val)
		if err != nil {
			return nil, err
		}

		paths := make([]battle_dto.BattleUnitPath, len(units))
		for idx, val := range units {
			paths[idx] = battle_dto.BattleUnitPath{
				UnitCode:        val.UnitCode,
				UnitLevel:       val.UnitLevel,
				UnitName:        val.UnitName,
				Origin:          val.Origin,
				ImageTypeNumber: rand.IntN(val.ImageTypeCount),
				Offense:         val.Offense,
				Defense:         val.Defense,
				Technique:       val.Technique,
				Speed:           val.Speed,
				Spirit:          val.Spirit,
				ElementID1:      val.ElementID1,
				ElementID2:      val.ElementID2,
			}
		}

		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}

		results[index] = battle_dto.BattleUnit{
			BattleUnitID: id.String(),
			Level:        levels[index],
			Paths:        paths,
		}
	}
	return results, nil
}

func (s *battleServiceImpl) GetPlayerTeamUnits(ctx context.Context, accountID, playerTeamID string) ([]battle_dto.BattleUnit, error) {
	playerUnits, err := s.teamRepo.GetPlayerUnitByTeamID(ctx, accountID, playerTeamID)
	if err != nil {
		return nil, err
	}

	maps := make(map[string]bool, 3)
	results := make([]battle_dto.BattleUnit, 3)
	i := -1

	for _, val := range playerUnits {
		_, ok := maps[val.PlayerUnitID]
		if !ok {
			i++
			maps[val.PlayerUnitID] = true
			results[i] = battle_dto.BattleUnit{
				BattleUnitID: val.PlayerUnitID,
				Level:        val.PlayerUnitLevel,
				Paths:        make([]battle_dto.BattleUnitPath, 0),
			}
		}

		results[i].Paths = append(results[i].Paths, battle_dto.BattleUnitPath{
			UnitCode:        val.UnitCode,
			UnitName:        val.UnitName,
			UnitLevel:       val.UnitLevel,
			Origin:          val.Origin,
			ImageTypeNumber: val.ImageTypeNumber,
			Offense:         val.Offense,
			Defense:         val.Defense,
			Technique:       val.Technique,
			Speed:           val.Speed,
			Spirit:          val.Spirit,
			ElementID1:      val.ElementID1,
			ElementID2:      val.ElementID2,
		})
	}
	return results, nil
}
