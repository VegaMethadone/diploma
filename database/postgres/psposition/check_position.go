package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/entity/position"

	"github.com/google/uuid"
)

func CheckPosition(positionId uuid.UUID) (*position.Position, error) {
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	findQuery := `
        SELECT id, company_id, lvl, name
        FROM positions
        WHERE id = $1
    `
	var pos position.Position
	err = db.QueryRow(findQuery, positionId).Scan(
		&pos.Id,
		&pos.CompanyId,
		&pos.Lvl,
		&pos.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	return &pos, nil
}
