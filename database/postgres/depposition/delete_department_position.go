package depposition

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (p PostgresDepPosition) DeleteDepartmentPosition(
	ctx context.Context,
	sharedTx *sql.Tx,
	positionID uuid.UUID,
) error {
	if sharedTx == nil {
		return errors.New("transaction must be started before query")
	}

	query := `
        DELETE FROM department_positions
        WHERE id = $1
    `

	result, err := sharedTx.ExecContext(ctx, query, positionID)
	if err != nil {
		return fmt.Errorf("failed to delete department position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department position with ID %s not found", positionID)
	}

	return nil
}
