package repository

import (
	"context"
	"log"
	unit_dto "tcg_card_battler/web-api/internal/dto/unit"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitRepository interface {
	GetUnitByCode(ctx context.Context, unitCode string) (*unit_dto.Unit, error)
	GetAllUnitLevelPathByCode(ctx context.Context, unitCode string) ([]unit_dto.GetUnitNextLevelPathRS, error)
	GetUnitLevelPathByCode(ctx context.Context, unitCode string, targetUnitCode string) (*unit_dto.GetUnitNextLevelPathRS, error)
	GetRandomUnitByLevel(ctx context.Context, level int) ([]unit_dto.Unit, error)
}

type unitRepositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewUnitRepository(pool *pgxpool.Pool) UnitRepository {
	return &unitRepositoryImpl{Pool: pool}
}

func (r *unitRepositoryImpl) GetUnitByCode(ctx context.Context, unitCode string) (*unit_dto.Unit, error) {
	query := `
		SELECT 
			unit_code, unit_name, origin, unit_level, offense, defense, technique, speed, spirit, tags, element_id_1, element_id_2
		FROM units
		WHERE unit_code = $1
		`

	var data unit_dto.Unit
	err := r.Pool.QueryRow(ctx, query, unitCode).Scan(
		&data.UnitCode,
		&data.UnitName,
		&data.Origin,
		&data.UnitLevel,
		&data.Offense,
		&data.Defense,
		&data.Technique,
		&data.Speed,
		&data.Spirit,
		&data.Tags,
		&data.ElementID1,
		&data.ElementID2,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *unitRepositoryImpl) GetAllUnitLevelPathByCode(ctx context.Context, unitCode string) ([]unit_dto.GetUnitNextLevelPathRS, error) {
	query := `
		SELECT 
			to_unit_code, u.unit_name, u.origin, u.tags, u.offense, u.defense, u.technique, u.speed, u.spirit, u.element_id_1, u.element_id_2
		FROM unit_level_paths ulp
		JOIN units u on ulp.to_unit_code = u.unit_code
		WHERE from_unit_code = $1`

	rows, err := r.Pool.Query(ctx, query, unitCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]unit_dto.GetUnitNextLevelPathRS, 0)
	for rows.Next() {
		var row unit_dto.GetUnitNextLevelPathRS

		err := rows.Scan(&row.UnitCode, &row.UnitName, &row.Origin, &row.Tags,
			&row.Offense, &row.Defense, &row.Technique, &row.Speed, &row.Spirit,
			&row.ElementID1, &row.ElementID2)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, row)
	}

	return results, nil
}

func (r *unitRepositoryImpl) GetUnitLevelPathByCode(ctx context.Context, unitCode string, targetUnitCode string) (*unit_dto.GetUnitNextLevelPathRS, error) {
	query := `
		SELECT 
			to_unit_code, u.unit_name, u.origin, u.unit_level, u.tags, u.offense, u.defense, u.technique, u.speed, u.spirit, u.element_id_1, u.element_id_2
		FROM unit_level_paths ulp
		JOIN units u on ulp.to_unit_code = u.unit_code
		WHERE from_unit_code = $1
		AND to_unit_code = $2`

	var data unit_dto.GetUnitNextLevelPathRS
	err := r.Pool.QueryRow(ctx, query, unitCode, targetUnitCode).Scan(
		&data.UnitCode,
		&data.UnitName,
		&data.Origin,
		&data.UnitLevel,
		&data.Tags,
		&data.Offense, &data.Defense, &data.Technique, &data.Speed, &data.Spirit,
		&data.ElementID1, &data.ElementID2,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *unitRepositoryImpl) GetRandomUnitByLevel(ctx context.Context, level int) ([]unit_dto.Unit, error) {
	query := `
		WITH RECURSIVE unit_path AS (
		    (SELECT 
		        u.unit_code, 
		        1 as depth
		    FROM units u
		    WHERE u.unit_level = $1
		    ORDER BY random()
		    LIMIT 1)

		    UNION ALL

		    SELECT 
		        child.from_unit_code as unit_code, 
		        up.depth + 1
		    FROM unit_path up
		    CROSS JOIN LATERAL (
		        SELECT next_ulp.from_unit_code, next_ulp.to_unit_code
		        FROM unit_level_paths next_ulp
		        WHERE next_ulp.to_unit_code = up.unit_code
		        ORDER BY random()
		        LIMIT 1
		    ) AS child
		)
		SELECT 
			u.unit_code, unit_name, origin, unit_level, u.offense, u.defense, u.technique, u.speed, u.spirit, image_type_count, tags, u.element_id_1, u.element_id_2
		FROM unit_path up
		JOIN units u on up.unit_code = u.unit_code
		ORDER BY unit_level;
	`
	rows, err := r.Pool.Query(ctx, query, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]unit_dto.Unit, 0)
	for rows.Next() {
		var row unit_dto.Unit

		err := rows.Scan(&row.UnitCode, &row.UnitName, &row.Origin, &row.UnitLevel,
			&row.Offense, &row.Defense, &row.Technique, &row.Speed, &row.Spirit, &row.ImageTypeCount, &row.Tags, &row.ElementID1, &row.ElementID2)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, row)
	}

	return results, nil
}
