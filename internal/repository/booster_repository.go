package repository

import (
	"context"
	booster_dto "tcg_card_battler/web-api/internal/dto/booster"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoosterRepository interface {
	GetAllBooster(ctx context.Context) (*booster_dto.GetAllBoosterRS, error)
	GetBoosterByCode(ctx context.Context, boosterCode string) (*booster_dto.Booster, error)
	GetHighestBoosterCard(ctx context.Context, boosterCode string, limit int) (*booster_dto.GetAllBoosterCardRS, error)
	GetAllBoosterCard(ctx context.Context, boosterCode string) (*booster_dto.GetAllBoosterCardRS, error)
	GetBoosterRarityRate(ctx context.Context, boosterCode string) (*booster_dto.GetBoosterRarityRateRS, error)
}

type boosterRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewBoosterRepository(pool *pgxpool.Pool) BoosterRepository {
	return &boosterRepositoryImpl{pool: pool}
}

func (s *boosterRepositoryImpl) GetAllBooster(ctx context.Context) (*booster_dto.GetAllBoosterRS, error) {
	query := `SELECT booster_code, booster_name, price FROM boosters`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Use pgx.CollectRows to eliminate manual loops and scanning errors
	boosters, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (booster_dto.Booster, error) {
		var b booster_dto.Booster
		// Scan into addresses of struct fields
		err := row.Scan(&b.BoosterCode, &b.BoosterName, &b.Price)
		return b, err
	})
	if err != nil {
		return nil, err
	}

	return &booster_dto.GetAllBoosterRS{Boosters: boosters}, nil
}

func (s *boosterRepositoryImpl) GetBoosterByCode(ctx context.Context, boosterCode string) (*booster_dto.Booster, error) {
	query := `SELECT booster_code, booster_name, price FROM boosters WHERE booster_code = $1`

	rows, err := s.pool.Query(ctx, query, boosterCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var b booster_dto.Booster
	err = s.pool.QueryRow(ctx, query, boosterCode).Scan(&b.BoosterCode, &b.BoosterName, &b.Price)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *boosterRepositoryImpl) GetHighestBoosterCard(ctx context.Context, boosterCode string, limit int) (*booster_dto.GetAllBoosterCardRS, error) {
	query := `
		SELECT 
		    c.card_code, 
		    c.image_type_number, 
		    c.card_rarity_code, 
			cr.price,
		    u.origin,
			u.element_id_1,
			u.element_id_2
		FROM booster_cards bc
		JOIN cards c ON bc.card_code = c.card_code AND bc.image_type_number = c.image_type_number
		JOIN card_rarities cr ON c.card_rarity_code = cr.card_rarity_code
		JOIN units u ON u.unit_code = c.card_code
		WHERE bc.booster_code = $1
		ORDER BY cr.price DESC
		LIMIT $2
	`

	rows, err := s.pool.Query(ctx, query, boosterCode, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (booster_dto.BoosterCard, error) {
		var c booster_dto.BoosterCard
		err := row.Scan(&c.CardCode, &c.ImageTypeNumber, &c.CardRarityCode, &c.Price, &c.Origin, &c.ElementID1, &c.ElementID2)
		return c, err
	})

	if err != nil {
		return nil, err
	}

	return &booster_dto.GetAllBoosterCardRS{Cards: cards}, nil
}

func (s *boosterRepositoryImpl) GetAllBoosterCard(ctx context.Context, boosterCode string) (*booster_dto.GetAllBoosterCardRS, error) {
	query := `
		SELECT 
		    c.card_code, 
		    c.image_type_number, 
		    c.card_rarity_code, 
			cr.price,
		    u.origin,
			u.element_id_1,
			u.element_id_2
		FROM booster_cards bc
		JOIN cards c ON bc.card_code = c.card_code AND bc.image_type_number = c.image_type_number
		JOIN card_rarities cr ON c.card_rarity_code = cr.card_rarity_code
		JOIN units u ON u.unit_code = c.card_code
		WHERE bc.booster_code = $1
		ORDER BY cr.price, c.card_code ASC
	`

	rows, err := s.pool.Query(ctx, query, boosterCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (booster_dto.BoosterCard, error) {
		var c booster_dto.BoosterCard
		err := row.Scan(&c.CardCode, &c.ImageTypeNumber, &c.CardRarityCode, &c.Price, &c.Origin, &c.ElementID1, &c.ElementID2)
		return c, err
	})

	if err != nil {
		return nil, err
	}

	return &booster_dto.GetAllBoosterCardRS{Cards: cards}, nil
}

func (s *boosterRepositoryImpl) GetBoosterRarityRate(ctx context.Context, boosterCode string) (*booster_dto.GetBoosterRarityRateRS, error) {
	query := `
		SELECT 
			bp.card_rarity_code, bp.percentage
		FROM booster_card_percentages bp
		WHERE bp.booster_code = $1
		ORDER BY bp.percentage ASC
	`

	rows, err := s.pool.Query(ctx, query, boosterCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (booster_dto.BoosterCardPercentage, error) {
		var c booster_dto.BoosterCardPercentage
		err := row.Scan(&c.CardRarityCode, &c.Percentage)
		return c, err
	})

	if err != nil {
		return nil, err
	}

	return &booster_dto.GetBoosterRarityRateRS{Items: cards}, nil
}
