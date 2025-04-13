package depemployee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (p PostgresEmployeeDepartment) ExistsEmployeeDepartment(
	ctx context.Context,
	sharedTx *sql.Tx,
	employeeID uuid.UUID,
	departmentID uuid.UUID,
) (bool, error) {
	if sharedTx == nil {
		return false, errors.New("transaction must be started before query")
	}

	query := `
        SELECT EXISTS(
            SELECT 1 
            FROM employee_department
            WHERE employee_id = $1
            AND department_id = $2
            AND is_active = true
        )
    `

	var exists bool
	err := sharedTx.QueryRowContext(ctx, query, employeeID, departmentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check employee department existence: %w", err)
	}

	return exists, nil
}
