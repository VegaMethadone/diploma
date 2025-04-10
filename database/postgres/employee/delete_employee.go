package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (p PostgresEmployee) DeleteEmployee(
	ctx context.Context,
	sharedTx *sql.Tx,
	employeeId uuid.UUID,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	query := `
        UPDATE employee_company
        SET 
            is_active = false,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		employeeId,
	)

	if err != nil {
		return fmt.Errorf("failed to deactivate employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found (id: %s)", employeeId)
	}

	return nil
}
