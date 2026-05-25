package repository

import (
	"context"
	general_dto "tcg_card_battler/web-api/internal/dto/general"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GeneralRepository interface {
	GetElements(ctx context.Context) ([]general_dto.Element, error)
	GetOrigins(ctx context.Context) ([]general_dto.Origin, error)
}

type generalRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewGeneralRepository(pool *pgxpool.Pool) GeneralRepository {
	return &generalRepositoryImpl{pool: pool}
}

func (r *generalRepositoryImpl) GetElements(ctx context.Context) ([]general_dto.Element, error) {
	query := `
		SELECT 
			element_id,
			element_name
		FROM elements
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	results := make([]general_dto.Element, 0)

	for rows.Next() {
		var row general_dto.Element

		err := rows.Scan(&row.ID, &row.Name)

		if err != nil {
			return nil, err
		}

		results = append(results, row)
	}

	return results, nil
}

func (r *generalRepositoryImpl) GetOrigins(ctx context.Context) ([]general_dto.Origin, error) {
	query := `
		SELECT 
			origin_code,
			origin_name
		FROM origins
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	results := make([]general_dto.Origin, 0)

	for rows.Next() {
		var row general_dto.Origin

		err := rows.Scan(&row.Code, &row.Name)

		if err != nil {
			return nil, err
		}

		results = append(results, row)
	}

	return results, nil
}
