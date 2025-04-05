package postgres

import (
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/entity/position"
)

func ChangePosition(what *position.Position) error {
	if what == nil {
		return fmt.Errorf("nil position provided")
	}

	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	updateQuery := `
        UPDATE positions
        SET company_id = $1,
            lvl = $2,
            name = $3
        WHERE id = $4
    `

	result, err := db.Exec(updateQuery,
		what.CompanyId,
		what.Lvl,
		what.Name,
		what.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("position with id %s not found", what.Id)
	}

	return nil
}
