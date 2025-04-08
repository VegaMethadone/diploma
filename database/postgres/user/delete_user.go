package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func DeleteUser(ctx context.Context, sharedTx *sql.Tx, id uuid.UUID) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}
	query := `
        DELETE FROM users
        WHERE id = $1
    `

	result, err := sharedTx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found (id: %s)", id)
	}

	return nil
}
