package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	inv_dto "tcg_card_battler/web-api/internal/dto/inventory"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository interface {
	exec(ctx context.Context) DBQuerier
	GetPlayerUnitCount(ctx context.Context, accountID string) (int, error)
	GetPlayerUnits(ctx context.Context, accountID string, limit, offset int) ([]inv_dto.PlayerUnit, error)
	GetPlayerUnitByIDs(ctx context.Context, accountID string, playerUnitIDs []string) ([]inv_dto.PlayerUnit, error)
	InvGetPlayerUnitDetailByID(ctx context.Context, accountID string, playerUnitID string) (*inv_dto.PlayerUnitDetailRS, error)
	InvGetPlayerUnitCardByUnitCode(ctx context.Context, accountID string, unitCode string) ([]inv_dto.PlayerCard, error)
	GetPlayerCards(ctx context.Context, accountID, unitCode string) (map[int]int, error)
	GetPlayerCardByCodeAndTypeNumber(ctx context.Context, accountID, cardCode string, typeNumbers []int) (map[int]int, error)
	GetAllPlayerCards(ctx context.Context, accountID string, limit int, cursorPrice int, cursorCardCode string, cursorImageTypeNumber int, pageNumber int, isPrev bool) (*inv_dto.GetAllPlayerCardRS, error)
	BatchInsertPlayerCards(ctx context.Context, accountID string, codes []string, imageTypes []int32, quantities []int32) error
	BatchUpdatePlayerCards(ctx context.Context, accountID, unitCode string, imageNums []int, qtys []int) error
	IncrementUnitLevel(ctx context.Context, playerUnitID string) error
	InvGetPlayerUnitPrevLevel(ctx context.Context, accountID string, playerUnitID string) ([]inv_dto.PlayerUnitPrevLevelRS, error)
	InvPostPlayerUnitChangeImage(ctx context.Context, accountID string, rq inv_dto.PlayerUnitLevelChangeImageRQ) error
	DecrementCard(ctx context.Context, accountID string, unitCode string, imageTypeNumber int) (int, error)
	DeleteCard(ctx context.Context, accountID string, rq inv_dto.PlayerUnitLevelChangeImageRQ) error
	BatchDeletePlayerCard(ctx context.Context, accountID, unitCode string, imageTypeNumbers []int) error
	UpdateUnitLevelImage(ctx context.Context, accountID string, rq inv_dto.PlayerUnitLevelChangeImageRQ) error
	InsertPlayerLevel(ctx context.Context, playerUnitID string, targetLevel int, unitCode string) error
	GetEligibleUnitsCount(ctx context.Context, accountID string) (int, error)
	GetEligibleUnitsList(ctx context.Context, accountID string, limit, offset int) ([]inv_dto.EligibleUnit, error)
	InsertPlayerUnit(ctx context.Context, accountID string) (string, error)
}

type inventoryRepositoryImpl struct {
	pool *pgxpool.Pool
}

func (r *inventoryRepositoryImpl) exec(ctx context.Context) DBQuerier {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}

	return r.pool
}

func NewInventoryRepository(pool *pgxpool.Pool) InventoryRepository {
	return &inventoryRepositoryImpl{pool: pool}
}

func (r *inventoryRepositoryImpl) GetPlayerUnitCount(ctx context.Context, accountID string) (int, error) {
	query := `
        SELECT 
			 COUNT(*)
        FROM player_units pu
		WHERE pu.account_id = $1`

	var total int
	err := r.pool.QueryRow(ctx, query, accountID).Scan(&total)
	return total, err
}

func (r *inventoryRepositoryImpl) GetPlayerUnits(ctx context.Context, accountID string, limit, offset int) ([]inv_dto.PlayerUnit, error) {
	query := `
        SELECT DISTINCT ON (pu.player_unit_id)
			 pu.player_unit_id, pu.player_unit_level, pul.unit_code, u.origin, pul.image_type_number, u.element_id_1, u.element_id_2
        FROM player_units pu
        JOIN player_unit_levels pul ON pul.player_unit_id = pu.player_unit_id
		JOIN units u ON pul.unit_code = u.unit_code
        WHERE pu.account_id = $1
		ORDER BY pu.player_unit_id DESC, pul.target_level DESC
		OFFSET $2
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, accountID, offset, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	results := make([]inv_dto.PlayerUnit, 0)

	for rows.Next() {
		var row inv_dto.PlayerUnit

		err := rows.Scan(&row.PlayerUnitID, &row.Level, &row.UnitCode, &row.Origin, &row.ImageTypeNumber, &row.ElementID1, &row.ElementID2)

		if err != nil {
			log.Fatal(err)
		}

		results = append(results, row)
	}

	return results, nil
}

func (r *inventoryRepositoryImpl) GetPlayerUnitByIDs(ctx context.Context, accountID string, playerUnitIDs []string) ([]inv_dto.PlayerUnit, error) {
	query := `
        SELECT DISTINCT ON (pu.player_unit_id)
			 pu.player_unit_id, pu.player_unit_level, pul.unit_code, u.origin, pul.image_type_number
        FROM player_units pu
        JOIN player_unit_levels pul ON pul.player_unit_id = pu.player_unit_id
		JOIN units u ON pul.unit_code = u.unit_code
        WHERE pu.account_id = $1
		AND pu.player_unit_id = ANY($2::uuid[]) 
		ORDER BY pu.player_unit_id DESC, pul.target_level DESC`

	rows, err := r.pool.Query(ctx, query, accountID, playerUnitIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]inv_dto.PlayerUnit, 0)

	for rows.Next() {
		var t inv_dto.PlayerUnit
		err := rows.Scan(&t.PlayerUnitID, &t.PlayerUnitID, &t.UnitCode, &t.Origin, &t.ImageTypeNumber)

		if err != nil {
			log.Fatal(err)
		}

		results = append(results, t)
	}

	return results, nil
}

func (r *inventoryRepositoryImpl) InvGetPlayerUnitDetailByID(ctx context.Context, accountID string, playerUnitID string) (*inv_dto.PlayerUnitDetailRS, error) {
	query := `
        SELECT DISTINCT ON (pu.player_unit_id)
		    pu.player_unit_id,
		    FIRST_VALUE(pul.unit_code) OVER (
		        PARTITION BY pu.player_unit_id 
		        ORDER BY pul.target_level ASC
		    ) as first_unit_code,
		    FIRST_VALUE(pul.unit_code) OVER (
		        PARTITION BY pu.player_unit_id 
		        ORDER BY pul.target_level DESC
		    ) as last_unit_code,
		    u.origin, 
		    u.unit_name, 
		    pu.player_unit_level, 
		    pul.image_type_number
		FROM player_units pu
		JOIN player_unit_levels pul ON pul.player_unit_id = pu.player_unit_id
		JOIN units u ON pul.unit_code = u.unit_code
		WHERE pu.account_id = $1
		  AND pu.player_unit_id = $2
		ORDER BY pu.player_unit_id, pul.target_level DESC;`

	var detail inv_dto.PlayerUnitDetailRS
	err := r.pool.QueryRow(ctx, query, accountID, playerUnitID).Scan(
		&detail.PlayerUnitID,
		&detail.FirstUnitCode,
		&detail.LastUnitCode,
		&detail.Origin,
		&detail.UnitName,
		&detail.PlayerUnitLevel,
		&detail.ImageTypeNumber,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &detail, nil
}

func (r *inventoryRepositoryImpl) InvGetPlayerUnitCardByUnitCode(ctx context.Context, accountID string, unitCode string) ([]inv_dto.PlayerCard, error) {
	query := `
        SELECT 
		 	pc.card_code, u.origin, image_type_number, qty
		FROM player_cards pc
		JOIN units u ON pc.card_code = u.unit_code
		WHERE 1=1
		AND pc.account_id = $1
		AND pc.card_code = $2`

	rows, err := r.pool.Query(ctx, query, accountID, unitCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]inv_dto.PlayerCard, 0)
	for rows.Next() {
		var row inv_dto.PlayerCard

		err := rows.Scan(&row.CardCode, &row.Origin, &row.ImageTypeNumber, &row.QTY)

		if err != nil {
			log.Fatal(err)
		}

		results = append(results, row)
	}

	return results, nil
}

func (r *inventoryRepositoryImpl) GetAllPlayerCards(ctx context.Context, accountID string, limit int, cursorPrice int, cursorCardCode string, cursorImageTypeNumber int, pageNumber int, isPrev bool) (*inv_dto.GetAllPlayerCardRS, error) {
	var totalRecords int
	query := `
		SELECT 
		    COUNT(*)
		FROM player_cards pc
		JOIN cards c on pc.card_code = c.card_code
		JOIN card_rarities cr on c.card_rarity_code = cr.card_rarity_code
		LEFT JOIN units u on c.card_code = u.unit_code
		WHERE pc.account_id = $1`
	err := r.pool.QueryRow(ctx, query, accountID).Scan(&totalRecords)
	if err != nil {
		return nil, err
	}
	result := &inv_dto.GetAllPlayerCardRS{
		TotalPage: 0,
		CurrPage:  pageNumber,
		Cards:     make([]inv_dto.PlayerCard, 0),
	}

	if limit > 0 {
		result.TotalPage = (totalRecords + limit - 1) / limit
	} else {
		result.TotalPage = 1
	}

	query = `
		SELECT 
		    c.card_code, c.image_type_number, u.origin, pc.qty, cr.card_rarity_code, cr.price
		FROM player_cards pc
		JOIN cards c on pc.card_code = c.card_code
		JOIN card_rarities cr on c.card_rarity_code = cr.card_rarity_code
		LEFT JOIN units u on c.card_code = u.unit_code
		WHERE pc.account_id = $1`

	args := []interface{}{accountID, limit}

	// 2. Add cursor logic only if not on page 1
	sortOrder := ` ORDER BY cr.price DESC, c.card_code ASC, c.image_type_number ASC`

	if pageNumber != 1 {
		var comparison string
		if isPrev {
			comparison = ` AND ((cr.price > $3) OR 
                      (cr.price = $3 AND c.card_code < $4) OR 
                      (cr.price = $3 AND c.card_code = $4 AND c.image_type_number < $5))`
			sortOrder = ` ORDER BY cr.price ASC, c.card_code DESC, c.image_type_number DESC`
		} else {
			comparison = ` AND ((cr.price < $3) OR 
                      (cr.price = $3 AND c.card_code > $4) OR 
                      (cr.price = $3 AND c.card_code = $4 AND c.image_type_number > $5))`
			sortOrder = ` ORDER BY cr.price DESC, c.card_code ASC, c.image_type_number ASC`
		}
		query += comparison
		args = append(args, cursorPrice, cursorCardCode, cursorImageTypeNumber)
	}

	// 3. Finalize string
	query += sortOrder + ` LIMIT $2`

	// 4. Execute with expanded args
	rows, err := r.pool.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row inv_dto.PlayerCard
		if err := rows.Scan(&row.CardCode, &row.ImageTypeNumber, &row.Origin, &row.QTY, &row.CardRarityCode, &row.Price); err != nil {
			return nil, err
		}
		result.Cards = append(result.Cards, row)
	}
	return result, nil
}

func (r *inventoryRepositoryImpl) GetPlayerCards(ctx context.Context, accountID, unitCode string) (map[int]int, error) {
	query := `SELECT image_type_number, qty FROM player_cards WHERE account_id = $1 AND card_code = $2`
	rows, err := r.pool.Query(ctx, query, accountID, unitCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make(map[int]int)
	for rows.Next() {
		var itn, q int
		if err := rows.Scan(&itn, &q); err != nil {
			return nil, err
		}
		cards[itn] = q
	}
	return cards, nil
}

func (r *inventoryRepositoryImpl) GetPlayerCardByCodeAndTypeNumber(ctx context.Context, accountID, cardCode string, typeNumbers []int) (map[int]int, error) {
	query := `
	SELECT 
		image_type_number, qty 
	FROM player_cards 
	WHERE account_id = $1 
	AND card_code = $2
	AND image_type_number = ANY($3)`

	rows, err := r.pool.Query(ctx, query, accountID, cardCode, typeNumbers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make(map[int]int)
	for rows.Next() {
		var itn, q int
		if err := rows.Scan(&itn, &q); err != nil {
			return nil, err
		}
		cards[itn] = q
	}
	return cards, nil
}

func (r *inventoryRepositoryImpl) BatchInsertPlayerCards(ctx context.Context, accountID string, codes []string, imageTypes []int32, quantities []int32) error {
	query := `
		INSERT INTO player_cards (account_id, card_code, image_type_number, qty)
		SELECT $1, unnest($2::text[]), unnest($3::int[]), unnest($4::int[])
		ON CONFLICT (account_id, card_code, image_type_number) DO UPDATE SET qty = player_cards.qty + EXCLUDED.qty;
	`
	_, err := r.exec(ctx).Exec(ctx, query, accountID, codes, imageTypes, quantities)
	return err
}

func (r *inventoryRepositoryImpl) BatchUpdatePlayerCards(ctx context.Context, accountID, unitCode string, imageNums []int, qtys []int) error {
	query := `
        UPDATE player_cards AS target SET qty = data.new_qty 
        FROM unnest($1::int[], $2::int[]) AS data(itn, new_qty)
        WHERE target.image_type_number = data.itn AND target.account_id = $3 AND target.card_code = $4`

	_, err := r.exec(ctx).Exec(ctx, query, imageNums, qtys, accountID, unitCode)
	return err
}

func (r *inventoryRepositoryImpl) IncrementUnitLevel(ctx context.Context, playerUnitID string) error {
	query := `UPDATE player_units SET player_unit_level = player_unit_level + 1 WHERE player_unit_id = $1`

	_, err := r.exec(ctx).Exec(ctx, query, playerUnitID)
	return err
}

func (r *inventoryRepositoryImpl) InvGetPlayerUnitPrevLevel(ctx context.Context, accountID string, playerUnitID string) ([]inv_dto.PlayerUnitPrevLevelRS, error) {
	query := `
		SELECT 
			pu.player_unit_id, pul.target_level, u.image_type_count, u.unit_name, u.unit_code, u.origin, pul.image_type_number, 
			e1.offense + e2.offense as offense,
			e1.defense + e2.defense as defense,
			e1.technique + e2.technique as technique,
			e1.speed + e2.speed as speed,
			e1.spirit + e2.spirit as spirit,
			u.element_id_1, u.element_id_2
		FROM player_units pu
		JOIN player_unit_levels pul on pu.player_unit_id = pul.player_unit_id
		JOIN units u on pul.unit_code = u.unit_code
		JOIN elements e1 on u.element_id_1 = e1.element_id
		JOIN elements e2 on u.element_id_2 = e2.element_id
		WHERE 1=1
		AND pu.account_id = $1
		AND pu.player_unit_id = $2
		ORDER BY pul.target_level ASC`

	rows, err := r.pool.Query(ctx, query, accountID, playerUnitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]inv_dto.PlayerUnitPrevLevelRS, 0)
	for rows.Next() {
		var row inv_dto.PlayerUnitPrevLevelRS

		err := rows.Scan(&row.PlayerUnitID, &row.TargetLevel, &row.ImageTypeCount, &row.UnitName, &row.UnitCode, &row.Origin, &row.ImageTypeNumber,
			&row.Offense, &row.Defense, &row.Technique, &row.Speed, &row.Spirit,
			&row.ElementID1, &row.ElementID2)

		if err != nil {
			log.Fatal(err)
		}

		results = append(results, row)
	}

	return results, nil
}

func (r *inventoryRepositoryImpl) InvPostPlayerUnitChangeImage(ctx context.Context, accountID string, rq inv_dto.PlayerUnitLevelChangeImageRQ) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Decrement QTY and get the new value in one go (Atomic)
	var newQty int
	updateCardQuery := `
		UPDATE player_cards 
		SET qty = qty - 1 
		WHERE account_id = $1 AND card_code = $2 AND image_type_number = $3 AND qty > 0
		RETURNING qty`

	err = tx.QueryRow(ctx, updateCardQuery, accountID, rq.UnitCode, rq.ImageTypeNumber).Scan(&newQty)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("not enough cards or card not found")
		}
		return err
	}

	// 2. Clean up if qty reached zero
	if newQty == 0 {
		deleteQuery := `DELETE FROM player_cards WHERE account_id = $1 AND card_code = $2 AND image_type_number = $3`
		if _, err = tx.Exec(ctx, deleteQuery, accountID, rq.UnitCode, rq.ImageTypeNumber); err != nil {
			return err
		}
	}

	// 3. Correct Postgres Update-From syntax
	updateLevelQuery := `
		UPDATE player_unit_levels pul
		SET image_type_number = $1
		FROM player_units pu
		WHERE pul.player_unit_id = pu.player_unit_id
		AND pu.account_id = $2
		AND pu.player_unit_id = $3
		AND pul.target_level = $4`

	result, err := tx.Exec(ctx, updateLevelQuery, rq.ImageTypeNumber, accountID, rq.PlayerUnitID, rq.TargetLevel)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("target level record not found")
	}

	return tx.Commit(ctx)
}

func (r *inventoryRepositoryImpl) DecrementCard(ctx context.Context, accountID string, unitCode string, imageTypeNumber int) (int, error) {
	var newQty int
	query := `UPDATE player_cards SET qty = qty - 1 WHERE account_id = $1 ... RETURNING qty`
	err := r.exec(ctx).QueryRow(ctx, query, accountID, unitCode, imageTypeNumber).Scan(&newQty)
	return newQty, err
}

func (r *inventoryRepositoryImpl) DeleteCard(ctx context.Context, accountID string, rq inv_dto.PlayerUnitLevelChangeImageRQ) error {
	query := `DELETE FROM player_cards WHERE account_id = $1 AND card_code = $2 AND image_type_number = $3`
	_, err := r.exec(ctx).Exec(ctx, query, accountID, rq.UnitCode, rq.ImageTypeNumber)
	return err
}

func (r *inventoryRepositoryImpl) BatchDeletePlayerCard(ctx context.Context, accountID, unitCode string, imageTypeNumbers []int) error {
	query := `DELETE FROM player_cards WHERE account_id = $1 AND card_code = $2 AND image_type_number = ANY($3)`
	_, err := r.exec(ctx).Exec(ctx, query, accountID, unitCode, imageTypeNumbers)
	return err
}

func (r *inventoryRepositoryImpl) UpdateUnitLevelImage(ctx context.Context, accountID string, rq inv_dto.PlayerUnitLevelChangeImageRQ) error {
	query := `UPDATE player_unit_levels pul SET image_type_number = $1
              FROM player_units pu WHERE pul.player_unit_id = pu.player_unit_id
              AND pu.account_id = $2 AND pu.player_unit_id = $3 AND pul.target_level = $4`
	_, err := r.exec(ctx).Exec(ctx, query, rq.ImageTypeNumber, accountID, rq.PlayerUnitID, rq.TargetLevel)
	return err
}

func (r *inventoryRepositoryImpl) InsertPlayerLevel(ctx context.Context, playerUnitID string, targetLevel int, unitCode string) error {
	query := `
		INSERT INTO player_unit_levels(player_unit_id, target_level, unit_code, image_type_number)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.exec(ctx).Exec(ctx, query, playerUnitID, targetLevel, unitCode, 0)
	return err
}

func (r *inventoryRepositoryImpl) GetEligibleUnitsCount(ctx context.Context, accountID string) (int, error) {
	// We wrap the GROUP BY query in a subquery to get the total count of groups
	query := `
        SELECT COUNT(*) FROM (
            SELECT 1
            FROM player_cards pc
            JOIN units u ON pc.card_code = u.unit_code
            WHERE pc.account_id = $1 AND u.unit_level = 1
            GROUP BY u.unit_code, u.origin
            HAVING sum(pc.qty) > 50
        ) AS groups`

	var total int
	err := r.pool.QueryRow(ctx, query, accountID).Scan(&total)
	return total, err
}

func (r *inventoryRepositoryImpl) GetEligibleUnitsList(ctx context.Context, accountID string, limit, offset int) ([]inv_dto.EligibleUnit, error) {
	query := `
        SELECT u.unit_code, u.origin
        FROM player_cards pc
        JOIN units u ON pc.card_code = u.unit_code
        WHERE pc.account_id = $1 AND u.unit_level = 1
        GROUP BY u.unit_code, u.origin
        HAVING sum(pc.qty) > 50
        ORDER BY u.unit_code, u.origin
        LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]inv_dto.EligibleUnit, 0)
	for rows.Next() {
		var row inv_dto.EligibleUnit

		err = rows.Scan(&row.UnitCode, &row.Origin)
		if err != nil {
			return nil, err
		}

		results = append(results, row)
	}
	return results, nil
}

func (r *inventoryRepositoryImpl) InsertPlayerUnit(ctx context.Context, accountID string) (string, error) {
	var playerUnitID string
	query := `
		INSERT INTO player_units (account_id, player_unit_level)
		VALUES ($1, 1)
		RETURNING player_unit_id;
	`
	err := r.exec(ctx).QueryRow(ctx, query, accountID).Scan(&playerUnitID)
	if err != nil {
		return "", err
	}
	return playerUnitID, nil
}
