package position

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (p PostgresPosition) DeletePosition(
	ctx context.Context,
	sharedTx *sql.Tx,
	positionId uuid.UUID,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	query := `
        UPDATE positions
        SET 
            is_active = false,
            updated_at = $1
        WHERE id = $2
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		time.Now(),
		positionId,
	)

	if err != nil {
		return fmt.Errorf("failed to delete position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("position not found (id: %s)", positionId)
	}

	return nil
}
