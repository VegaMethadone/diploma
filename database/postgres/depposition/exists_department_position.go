package depposition

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (p PostgresDepPosition) ExistsPosition(
	ctx context.Context,
	sharedTx *sql.Tx,
	positionID uuid.UUID,
) (bool, error) {
	if sharedTx == nil {
		return false, errors.New("transaction must be started before query")
	}

	query := `
        SELECT EXISTS(
            SELECT 1 
            FROM department_positions 
            WHERE id = $1
        )
    `

	var exists bool
	err := sharedTx.QueryRowContext(ctx, query, positionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check position existence: %w", err)
	}

	return exists, nil
}
