package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	inv_dto "tcg_card_battler/web-api/internal/dto/inventory"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository interface {
	exec(ctx context.Context) DBQuerier
	GetPlayerUnitCount(ctx context.Context, accountID, name, origin string, element1, element2, level, lastUnitLevel int) (int, error)
	GetPlayerUnits(ctx context.Context, accountID, name, origin string, limit, offset, element1, element2, level, lastUnitLevel, sort int) ([]inv_dto.PlayerUnit, error)
	GetPlayerUnitByIDs(ctx context.Context, accountID string, playerUnitIDs []string) ([]inv_dto.PlayerUnit, error)
	InvGetPlayerUnitDetailByID(ctx context.Context, accountID string, playerUnitID string) (*inv_dto.PlayerUnitDetailRS, error)
	InvGetPlayerUnitCardByUnitCode(ctx context.Context, accountID string, unitCode string) ([]inv_dto.PlayerCard, error)
	GetPlayerCards(ctx context.Context, accountID, unitCode string) (map[int]int, error)
	GetPlayerCardByCodeAndTypeNumber(ctx context.Context, accountID, cardCode string, typeNumbers []int) (map[int]int, error)
	GetPlayerCardCount(ctx context.Context, accountID string) (int, error)
	GetAllPlayerCards(ctx context.Context, accountID string, limit, offset int) ([]inv_dto.PlayerCard, error)
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
	InsertPlayerUnit(ctx context.Context, accountID string, level int) (string, error)
	GetPlayerUnitPriceByID(ctx context.Context, playerUnitID string) (int, error)
	SellPlayerUnit(ctx context.Context, accountID, playerunitID string) error
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

func (r *inventoryRepositoryImpl) GetPlayerUnitCount(ctx context.Context, accountID string, name, origin string, element1, element2, level, lastUnitLevel int) (int, error) {
	baseQuery := `
		FROM player_units pu
		CROSS JOIN LATERAL (
			SELECT unit_code, target_level
			FROM player_unit_levels
			WHERE player_unit_id = pu.player_unit_id
			ORDER BY target_level DESC
			LIMIT 1
		) pul
		JOIN units u ON pul.unit_code = u.unit_code
		WHERE pu.account_id = $1`

	args := []any{accountID}
	filterSQL, args := buildInventoryFilters(name, origin, element1, element2, level, lastUnitLevel, args)
	finalQuery := fmt.Sprintf("SELECT COUNT(*) %s %s", baseQuery, filterSQL)

	var total int
	err := r.pool.QueryRow(ctx, finalQuery, args...).Scan(&total)
	return total, err
}

func (r *inventoryRepositoryImpl) GetPlayerUnits(ctx context.Context, accountID, name, origin string, limit, offset, element1, element2, level, lastUnitLevel, sort int) ([]inv_dto.PlayerUnit, error) {
	// 1. Base Query: Use LATERAL to lock in exactly ONE latest level record per unit
	query := `
		SELECT 
			pu.player_unit_id, pu.player_unit_level, pul.unit_code, 
			u.origin, pul.image_type_number, u.element_id_1, u.element_id_2
		FROM player_units pu
		CROSS JOIN LATERAL (
			SELECT unit_code, image_type_number, target_level
			FROM player_unit_levels
			WHERE player_unit_id = pu.player_unit_id
			ORDER BY target_level DESC
			LIMIT 1
		) pul
		JOIN units u ON pul.unit_code = u.unit_code
		WHERE pu.account_id = $1`

	args := []any{accountID}

	// 2. Apply Filters (These now act on the LATEST record only)
	filterSQL, args := buildInventoryFilters(name, origin, element1, element2, level, lastUnitLevel, args)

	sortSQL := "ORDER BY pu.created_time ASC"
	switch sort {
	case 1:
		sortSQL = "ORDER BY pu.created_time DESC"
		break
	case 2:
		sortSQL = "ORDER BY pu.player_unit_level ASC"
		break
	case 3:
		sortSQL = "ORDER BY pu.player_unit_level DESC"
		break
	}

	finalQuery := fmt.Sprintf(`
		%s %s
		%s
		LIMIT $%d OFFSET $%d`,
		query, filterSQL, sortSQL, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, finalQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]inv_dto.PlayerUnit, 0)
	for rows.Next() {
		var row inv_dto.PlayerUnit
		if err := rows.Scan(&row.PlayerUnitID, &row.Level, &row.UnitCode, &row.Origin, &row.ImageTypeNumber, &row.ElementID1, &row.ElementID2); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

func buildInventoryFilters(name, origin string, e1, e2, lvl, lastLvl int, args []any) (string, []any) {
	var sql strings.Builder

	if name != "" {
		args = append(args, "%"+name+"%")
		fmt.Fprintf(&sql, " AND u.unit_name ILIKE $%d", len(args))
	}
	if origin != "" && origin != "000" {
		args = append(args, origin)
		fmt.Fprintf(&sql, " AND u.origin = $%d", len(args))
	}
	if lvl > 0 {
		args = append(args, lvl)
		fmt.Fprintf(&sql, " AND pu.player_unit_level = $%d", len(args))
	}
	if lastLvl > 0 {
		args = append(args, lastLvl)
		fmt.Fprintf(&sql, " AND pul.target_level = $%d", len(args))
	}

	if e1 > 0 && e2 > 0 {
		args = append(args, e1, e2)
		fmt.Fprintf(&sql, " AND ((u.element_id_1 = $%d AND u.element_id_2 = $%d)", len(args)-1, len(args))
		args = append(args, e2, e1)
		fmt.Fprintf(&sql, " OR (u.element_id_1 = $%d AND u.element_id_2 = $%d))", len(args)-1, len(args))
	} else if e1 > 0 || e2 > 0 {
		val := e1
		if e2 > 0 {
			val = e2
		}
		args = append(args, val)
		p := len(args)
		fmt.Fprintf(&sql, " AND (u.element_id_1 = $%d OR u.element_id_2 = $%d)", p, p)
	}

	return sql.String(), args
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
		WITH pul AS (
		    SELECT DISTINCT(player_unit_id)
		        player_unit_id,
				FIRST_VALUE(unit_code)	OVER (PARTITION BY player_unit_id ORDER BY target_level ASC) as first_unit_code,
		    	FIRST_VALUE(unit_code)	OVER (PARTITION BY player_unit_id ORDER BY target_level DESC) as last_unit_code,
		    	FIRST_VALUE(image_type_number)	OVER (PARTITION BY player_unit_id ORDER BY target_level DESC) as image_type_number
			FROM player_unit_levels
			GROUP BY player_unit_id, target_level
		)
        SELECT
		    pu.player_unit_id,
		    pul.first_unit_code,
		    pul.last_unit_code,
		    u.origin, 
		    u.unit_name, 
		    pu.player_unit_level, 
		    pul.image_type_number,
			u.element_id_1,
			u.element_id_2
		FROM player_units pu
		JOIN pul ON pu.player_unit_id = pul.player_unit_id
		JOIN units u ON pul.last_unit_code = u.unit_code
		WHERE pu.account_id = $1
		  AND pu.player_unit_id = $2
		ORDER BY pu.created_time ASC;`

	var detail inv_dto.PlayerUnitDetailRS
	err := r.pool.QueryRow(ctx, query, accountID, playerUnitID).Scan(
		&detail.PlayerUnitID,
		&detail.FirstUnitCode,
		&detail.LastUnitCode,
		&detail.Origin,
		&detail.UnitName,
		&detail.PlayerUnitLevel,
		&detail.ImageTypeNumber,
		&detail.ElementID1,
		&detail.ElementID2,
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
		 	pc.card_code, u.origin, pc.image_type_number, qty
		FROM player_cards pc
		JOIN cards c ON pc.card_code = c.card_code AND pc.image_type_number = c.image_type_number
		JOIN card_rarities cr ON c.card_rarity_code = cr.card_rarity_code
		JOIN units u ON pc.card_code = u.unit_code
		WHERE 1=1
		AND pc.account_id = $1
		AND pc.card_code = $2
		ORDER BY cr.price
		`

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

func (r *inventoryRepositoryImpl) GetPlayerCardCount(ctx context.Context, accountID string) (int, error) {
	query := `
        SELECT 
		    COUNT(*)
		FROM player_cards pc
		WHERE pc.account_id = $1`

	var total int
	err := r.pool.QueryRow(ctx, query, accountID).Scan(&total)
	return total, err
}

func (r *inventoryRepositoryImpl) GetAllPlayerCards(ctx context.Context, accountID string, limit, offset int) ([]inv_dto.PlayerCard, error) {
	query := `
		SELECT 
		    c.card_code, c.image_type_number, u.origin, pc.qty, cr.card_rarity_code, cr.price
		FROM player_cards pc
		JOIN cards c on pc.card_code = c.card_code AND pc.image_type_number = c.image_type_number
		JOIN card_rarities cr on c.card_rarity_code = cr.card_rarity_code
		JOIN units u on c.card_code = u.unit_code
		WHERE pc.account_id = $1
		ORDER BY cr.price DESC, c.card_code ASC, c.image_type_number ASC
		LIMIT $2
		OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, accountID, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]inv_dto.PlayerCard, 0)
	for rows.Next() {
		var row inv_dto.PlayerCard
		if err := rows.Scan(&row.CardCode, &row.ImageTypeNumber, &row.Origin, &row.QTY, &row.CardRarityCode, &row.Price); err != nil {
			return nil, err
		}
		cards = append(cards, row)
	}
	return cards, nil
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
	query := `
	UPDATE player_cards 
		SET qty = qty - 1 
	WHERE account_id = $1
	AND card_code = $2
	AND image_type_number = $3
	RETURNING qty`
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

func (r *inventoryRepositoryImpl) InsertPlayerUnit(ctx context.Context, accountID string, level int) (string, error) {
	var playerUnitID string
	query := `
		INSERT INTO player_units (account_id, player_unit_level)
		VALUES ($1, $2)
		RETURNING player_unit_id;
	`
	err := r.exec(ctx).QueryRow(ctx, query, accountID, level).Scan(&playerUnitID)
	if err != nil {
		return "", err
	}
	return playerUnitID, nil
}

func (r *inventoryRepositoryImpl) GetPlayerUnitPriceByID(ctx context.Context, playerUnitID string) (int, error) {
	gold := 0
	query := `
		SELECT 
			SUM(total) as total
		FROM(
			SELECT 
				SUM(cr.price) as total
			FROM player_unit_levels pul
			JOIN cards c ON pul.unit_code = c.card_code AND pul.image_type_number = c.image_type_number
			JOIN card_rarities cr ON c.card_rarity_code = cr.card_rarity_code
			WHERE player_unit_id = $1

			UNION ALL

			SELECT 
				SUM(cr.price) * 19 as total
			FROM player_unit_levels pul
			JOIN cards c ON pul.unit_code = c.card_code AND c.image_type_number = 0
			JOIN card_rarities cr ON c.card_rarity_code = cr.card_rarity_code
			WHERE player_unit_id = $1
			AND pul.target_level != 1

			UNION ALL

			SELECT 
				(player_unit_level - 1) * 10 + 49 as total
			FROM player_units pu
			WHERE player_unit_id = $1
		)x;
	`
	err := r.exec(ctx).QueryRow(ctx, query, playerUnitID).Scan(&gold)
	if err != nil {
		return 0, err
	}
	return gold, nil
}

func (r *inventoryRepositoryImpl) SellPlayerUnit(ctx context.Context, accountID, playerunitID string) error {
	query := `
		DELETE FROM player_units
		WHERE account_id = $1
		AND player_unit_id = $2
	`

	_, err := r.exec(ctx).Exec(ctx, query, accountID, playerunitID)
	return err
}
