package depemployee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (p PostgresEmployeeDepartment) DeleteEmployeeDepartment(
	ctx context.Context,
	sharedTx *sql.Tx,
	demployeeID uuid.UUID,
) error {
	if sharedTx == nil {
		return errors.New("transaction must be started before query")
	}

	query := `
        DELETE FROM employee_department
        WHERE id = $1
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		demployeeID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete employee department link: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee department link not found (id: %s)", demployeeID)
	}

	return nil
}
