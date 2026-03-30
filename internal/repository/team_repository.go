package repository

import (
	"context"
	team_dto "tcg_card_battler/web-api/internal/dto/team"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository interface {
	GetPlayerTeamCount(ctx context.Context, accountID string) (int, error)
	GetPlayerTeam(ctx context.Context, accountID string, limit, offset int) ([]team_dto.PlayerTeam, error)
	GetPlayerTeamByTeamID(ctx context.Context, accountID, teamID string) (*team_dto.PlayerTeam, error)
	InsertPlayerTeam(ctx context.Context, accountID, teamName, playerUnitID1, playerUnitID2, PlayerUnitID3 string) error
	GetActivePlayerTeamID(ctx context.Context, accountID string) (string, error)
	UnsetActivePlayerTeam(ctx context.Context, accountID string) error
	SetActivePlayerTeam(ctx context.Context, accountID, teamID string) error
	DeletePlayerTeam(ctx context.Context, accountID, playerTeamID string) error
	GetPlayerUnitByTeamID(ctx context.Context, accountID, playerTeamID string) ([]team_dto.PlayerTeamUnit, error)
}
type teamRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(p *pgxpool.Pool) TeamRepository {
	return &teamRepositoryImpl{pool: p}
}

func (r *teamRepositoryImpl) GetPlayerTeamCount(ctx context.Context, accountID string) (int, error) {
	totalRecords := 0
	query := `
		SELECT 
			COUNT(*)
		FROM player_teams pt
 		WHERE pt.account_id = $1
	`
	err := r.pool.QueryRow(ctx, query, accountID).Scan(&totalRecords)
	return totalRecords, err
}

func (r *teamRepositoryImpl) GetPlayerTeam(ctx context.Context, accountID string, limit, offset int) ([]team_dto.PlayerTeam, error) {
	query := `
		SELECT 
			pt.player_team_id,
			pt.team_name, 
			pt.is_active,
			(u1.player_unit_level + u2.player_unit_level + u3.player_unit_level) as TeamLevel,
			u1.*,
			u2.*,
			u3.*
		FROM player_teams pt
		LEFT JOIN LATERAL (
		    SELECT
				 pu.player_unit_id, pu.player_unit_level, pul.unit_code, u.origin, pul.image_type_number, u.element_id_1, u.element_id_2
		    FROM player_units pu
		    JOIN player_unit_levels pul ON pul.player_unit_id = pu.player_unit_id
			JOIN units u ON pul.unit_code = u.unit_code 
			WHERE pu.player_unit_id = pt.player_unit_id_1
		    ORDER BY pul.target_level DESC LIMIT 1
		) u1 ON true
		LEFT JOIN LATERAL (
		    SELECT
				 pu.player_unit_id, pu.player_unit_level, pul.unit_code, u.origin, pul.image_type_number, u.element_id_1, u.element_id_2
		    FROM player_units pu
		    JOIN player_unit_levels pul ON pul.player_unit_id = pu.player_unit_id
			JOIN units u ON pul.unit_code = u.unit_code 
			WHERE pu.player_unit_id = pt.player_unit_id_2
		    ORDER BY pul.target_level DESC LIMIT 1
		) u2 ON true
		LEFT JOIN LATERAL (
		    SELECT
				 pu.player_unit_id, pu.player_unit_level, pul.unit_code, u.origin, pul.image_type_number, u.element_id_1, u.element_id_2
		    FROM player_units pu
		    JOIN player_unit_levels pul ON pul.player_unit_id = pu.player_unit_id
			JOIN units u ON pul.unit_code = u.unit_code 
			WHERE pu.player_unit_id = pt.player_unit_id_3
		    ORDER BY pul.target_level DESC LIMIT 1
		) u3 ON true
 		WHERE pt.account_id = $1
		ORDER BY pt.is_active DESC, pt.player_team_id DESC
		OFFSET $2
		LIMIT $3
	`
	rows, err := r.pool.Query(ctx, query, accountID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]team_dto.PlayerTeam, 0)
	for rows.Next() {
		var t team_dto.PlayerTeam

		err := rows.Scan(&t.TeamID, &t.TeamName, &t.IsActive, &t.TeamLevel,
			&t.PlayerUnit1.PlayerUnitID, &t.PlayerUnit1.Level, &t.PlayerUnit1.UnitCode, &t.PlayerUnit1.Origin, &t.PlayerUnit1.ImageTypeNumber, &t.PlayerUnit1.ElementID1, &t.PlayerUnit1.ElementID2,
			&t.PlayerUnit2.PlayerUnitID, &t.PlayerUnit2.Level, &t.PlayerUnit2.UnitCode, &t.PlayerUnit2.Origin, &t.PlayerUnit2.ImageTypeNumber, &t.PlayerUnit2.ElementID1, &t.PlayerUnit2.ElementID2,
			&t.PlayerUnit3.PlayerUnitID, &t.PlayerUnit3.Level, &t.PlayerUnit3.UnitCode, &t.PlayerUnit3.Origin, &t.PlayerUnit3.ImageTypeNumber, &t.PlayerUnit3.ElementID1, &t.PlayerUnit3.ElementID2)

		if err != nil {
			return nil, err
		}

		results = append(results, t)
	}

	return results, nil
}

func (r *teamRepositoryImpl) GetPlayerTeamByTeamID(ctx context.Context, accountID, teamID string) (*team_dto.PlayerTeam, error) {
	query := `
		SELECT 
			pt.player_team_id,
			pt.team_name, 
			pt.is_active,
			(u1.player_unit_level + u2.player_unit_level + u3.player_unit_level) as TeamLevel,
			u1.*,
			u2.*,
			u3.*
		FROM player_teams pt
		LEFT JOIN LATERAL (
		    SELECT
				 pu.player_unit_id, pu.player_unit_level, pul.unit_code, u.origin, pul.image_type_number, u.element_id_1, u.element_id_2
		    FROM player_units pu
		    JOIN player_unit_levels pul ON pul.player_unit_id = pu.player_unit_id
			JOIN units u ON pul.unit_code = u.unit_code 
			WHERE pu.player_unit_id = pt.player_unit_id_1
		    ORDER BY pul.target_level DESC LIMIT 1
		) u1 ON true
		LEFT JOIN LATERAL (
		    SELECT
				 pu.player_unit_id, pu.player_unit_level, pul.unit_code, u.origin, pul.image_type_number, u.element_id_1, u.element_id_2
		    FROM player_units pu
		    JOIN player_unit_levels pul ON pul.player_unit_id = pu.player_unit_id
			JOIN units u ON pul.unit_code = u.unit_code 
			WHERE pu.player_unit_id = pt.player_unit_id_2
		    ORDER BY pul.target_level DESC LIMIT 1
		) u2 ON true
		LEFT JOIN LATERAL (
		    SELECT
				 pu.player_unit_id, pu.player_unit_level, pul.unit_code, u.origin, pul.image_type_number, u.element_id_1, u.element_id_2
		    FROM player_units pu
		    JOIN player_unit_levels pul ON pul.player_unit_id = pu.player_unit_id
			JOIN units u ON pul.unit_code = u.unit_code 
			WHERE pu.player_unit_id = pt.player_unit_id_3
		    ORDER BY pul.target_level DESC LIMIT 1
		) u3 ON true
 		WHERE pt.account_id = $1
		AND pt.player_team_id = $2
	`
	var t team_dto.PlayerTeam
	err := r.pool.QueryRow(ctx, query, accountID, teamID).Scan(&t.TeamID, &t.TeamName, &t.IsActive, &t.TeamLevel,
		&t.PlayerUnit1.PlayerUnitID, &t.PlayerUnit1.Level, &t.PlayerUnit1.UnitCode, &t.PlayerUnit1.Origin, &t.PlayerUnit1.ImageTypeNumber, &t.PlayerUnit1.ElementID1, &t.PlayerUnit1.ElementID2,
		&t.PlayerUnit2.PlayerUnitID, &t.PlayerUnit2.Level, &t.PlayerUnit2.UnitCode, &t.PlayerUnit2.Origin, &t.PlayerUnit2.ImageTypeNumber, &t.PlayerUnit2.ElementID1, &t.PlayerUnit2.ElementID2,
		&t.PlayerUnit3.PlayerUnitID, &t.PlayerUnit3.Level, &t.PlayerUnit3.UnitCode, &t.PlayerUnit3.Origin, &t.PlayerUnit3.ImageTypeNumber, &t.PlayerUnit3.ElementID1, &t.PlayerUnit3.ElementID2)

	return &t, err
}

func (r *teamRepositoryImpl) InsertPlayerTeam(ctx context.Context, accountID, teamName, playerUnitID1, playerUnitID2, PlayerUnitID3 string) error {
	query := `
		INSERT INTO player_teams(account_id, team_name, player_unit_id_1, player_unit_id_2, player_unit_id_3)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, query, accountID, teamName, playerUnitID1, playerUnitID2, PlayerUnitID3)
	if err != nil {
		return err
	}
	return nil
}

func (r *teamRepositoryImpl) UnsetActivePlayerTeam(ctx context.Context, accountID string) error {
	query := `
		UPDATE player_teams SET is_active = '0'
		WHERE is_active = '1'
		and account_id = $1`

	_, err := r.pool.Exec(ctx, query, accountID)
	if err != nil {
		return err
	}
	return nil
}

func (r *teamRepositoryImpl) GetActivePlayerTeamID(ctx context.Context, accountID string) (string, error) {
	teamID := ""
	query := `
		SELECT 
			player_team_id
		FROM player_teams pt
		WHERE pt.is_active = '1'
		AND pt.account_id = $1
		LIMIT 1
	`
	err := r.pool.QueryRow(ctx, query, accountID).Scan(&teamID)
	return teamID, err
}

func (r *teamRepositoryImpl) SetActivePlayerTeam(ctx context.Context, accountID, teamID string) error {
	query := `
		UPDATE player_teams SET is_active = '1'
		WHERE account_id = $1
		AND player_team_id = $2`

	_, err := r.pool.Exec(ctx, query, accountID, teamID)
	if err != nil {
		return err
	}
	return nil
}

func (r *teamRepositoryImpl) DeletePlayerTeam(ctx context.Context, accountID, playerTeamID string) error {
	query := `
		DELETE FROM player_teams
		WHERE account_id = $1
		AND player_team_id = $2`

	_, err := r.pool.Exec(ctx, query, accountID, playerTeamID)
	if err != nil {
		return nil
	}
	return nil
}

func (r *teamRepositoryImpl) GetPlayerUnitByTeamID(ctx context.Context, accountID, playerTeamID string) ([]team_dto.PlayerTeamUnit, error) {
	query := `
		SELECT 
			pu.player_unit_id, pu.player_unit_level, pul.image_type_number, u.unit_code, unit_name, origin, unit_level, tags,
			offense, defense, technique, speed, spirit, u.element_id_1, u.element_id_2
		FROM player_units pu
		JOIN player_teams pt ON pu.player_unit_id IN (pt.player_unit_id_1, pt.player_unit_id_2, pt.player_unit_id_3)
		JOIN player_unit_levels pul on pu.player_unit_id = pul.player_unit_id
		JOIN units u on pul.unit_code= u.unit_code
		WHERE pt.player_team_id = $2
		AND pu.account_id = $1
		ORDER BY 
		CASE 
			WHEN pu.player_unit_id = pt.player_unit_id_1 THEN 1
			WHEN pu.player_unit_id = pt.player_unit_id_2 THEN 2
			ELSE 3
		END, u.unit_level
	`

	rows, err := r.pool.Query(ctx, query, accountID, playerTeamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]team_dto.PlayerTeamUnit, 0)
	for rows.Next() {
		var r team_dto.PlayerTeamUnit

		err := rows.Scan(&r.PlayerUnitID, &r.PlayerUnitLevel, &r.ImageTypeNumber,
			&r.UnitCode, &r.UnitName, &r.Origin, &r.UnitLevel, &r.Tags,
			&r.Offense, &r.Defense, &r.Technique, &r.Speed, &r.Spirit, &r.ElementID1, &r.ElementID2,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}
