package department

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (p PostgresDepartment) DeleteDepartment(
	ctx context.Context,
	sharedTx *sql.Tx,
	departmentId uuid.UUID,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	query := `
        UPDATE departments
        SET 
            is_active = false,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
        AND is_active = true
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		departmentId,
	)

	if err != nil {
		return fmt.Errorf("failed to deactivate department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found or already inactive (id: %s)", departmentId)
	}

	return nil
}
