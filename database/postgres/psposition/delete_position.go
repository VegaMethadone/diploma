package postgres

import (
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"

	"github.com/google/uuid"
)

func DeletePosition(positionId uuid.UUID) error {
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	deleteQuery := `
        DELETE FROM positions
        WHERE id = $1
    `

	result, err := db.Exec(deleteQuery, positionId)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("position with id %s not found", positionId)
	}

	return nil
}
