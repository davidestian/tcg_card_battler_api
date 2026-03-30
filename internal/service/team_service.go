package service

import (
	"context"
	"fmt"
	team_dto "tcg_card_battler/web-api/internal/dto/team"
	"tcg_card_battler/web-api/internal/repository"
)

type TeamService interface {
	GetPlayerTeam(ctx context.Context, accountID string, limit, page int) (*team_dto.GetPlayerTeamRS, error)
	GetPlayerTeamByTeamID(ctx context.Context, accountID, teamID string) (*team_dto.PlayerTeam, error)
	GetActivePlayerTeamID(ctx context.Context, accountID string) (string, error)
	PostPlayerTeam(ctx context.Context, accountID string, rq team_dto.PostPlayerTeamRQ) error
	PutActivePlayerTeam(ctx context.Context, accountID, teamID string) error
	DeletePlayerTeam(ctx context.Context, accountID, playerTeamID string) error
}

type teamServiceImpl struct {
	teamRepo      repository.TeamRepository
	inventoryRepo repository.InventoryRepository
}

func NewTeamService(tr repository.TeamRepository, ir repository.InventoryRepository) TeamService {
	return &teamServiceImpl{teamRepo: tr, inventoryRepo: ir}
}

func (s *teamServiceImpl) GetPlayerTeam(ctx context.Context, accountID string, limit, page int) (*team_dto.GetPlayerTeamRS, error) {
	totalRecords, err := s.teamRepo.GetPlayerTeamCount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	teams, err := s.teamRepo.GetPlayerTeam(ctx, accountID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalPage := 0
	if limit > 0 {
		totalPage = (totalRecords + limit - 1) / limit
	}

	return &team_dto.GetPlayerTeamRS{
		TotalPage:   totalPage,
		PlayerTeams: teams,
	}, err
}

func (s *teamServiceImpl) GetPlayerTeamByTeamID(ctx context.Context, accountID, teamID string) (*team_dto.PlayerTeam, error) {
	team, err := s.teamRepo.GetPlayerTeamByTeamID(ctx, accountID, teamID)
	return team, err
}

func (s *teamServiceImpl) GetActivePlayerTeamID(ctx context.Context, accountID string) (string, error) {
	teamID, err := s.teamRepo.GetActivePlayerTeamID(ctx, accountID)
	return teamID, err
}

func (s *teamServiceImpl) PostPlayerTeam(ctx context.Context, accountID string, rq team_dto.PostPlayerTeamRQ) error {
	playerUnitIDs := make([]string, 3)
	playerUnitIDs[0] = rq.PlayerUnitID1
	playerUnitIDs[1] = rq.PlayerUnitID2
	playerUnitIDs[2] = rq.PlayerUnitID3

	playerUnits, err := s.inventoryRepo.GetPlayerUnitByIDs(ctx, accountID, playerUnitIDs)
	if err != nil {
		return err
	}
	if len(playerUnits) != 3 {
		return fmt.Errorf("Player Units not found")
	}
	if rq.TeamName == "" {
		rq.TeamName = "no name"
	}
	err = s.teamRepo.InsertPlayerTeam(ctx, accountID, rq.TeamName, rq.PlayerUnitID1, rq.PlayerUnitID2, rq.PlayerUnitID3)
	if err != nil {
		return err
	}
	return err
}

func (s *teamServiceImpl) PutActivePlayerTeam(ctx context.Context, accountID, teamID string) error {
	err := s.teamRepo.UnsetActivePlayerTeam(ctx, accountID)
	if err != nil {
		return err
	}
	err = s.teamRepo.SetActivePlayerTeam(ctx, accountID, teamID)
	if err != nil {
		return err
	}
	return err
}

func (s *teamServiceImpl) DeletePlayerTeam(ctx context.Context, accountID, playerTeamID string) error {
	teamID, err := s.teamRepo.GetActivePlayerTeamID(ctx, accountID)
	if err != nil {
		return err
	}
	if teamID == playerTeamID {
		return fmt.Errorf("Cannot delete active team")
	}

	err = s.teamRepo.DeletePlayerTeam(ctx, accountID, playerTeamID)
	if err != nil {
		return err
	}
	return err
}
