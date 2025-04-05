package postgres

import (
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/entity/position"

	"github.com/google/uuid"
)

func GetPositions(companyId uuid.UUID) ([]position.Position, error) {
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Исправлена опечатка в имени поля (company_id вместо compnay_id)
	getQuery := `
        SELECT id, company_id, lvl, name 
        FROM positions
        WHERE company_id = $1
    `

	rows, err := db.Query(getQuery, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var positions []position.Position

	for rows.Next() {
		var pos position.Position
		if err := rows.Scan(&pos.Id, &pos.CompanyId, &pos.Lvl, &pos.Name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		positions = append(positions, pos)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return positions, nil
}

/*
	merged my code again and pass CI except cangjie manifest
	got comments from dmitry and andrei
	andrei told me to rename some veriables, make them immutable and rewrite some code
	cause i made an error of the process a part of ResultContainer but it should be enum
*/
