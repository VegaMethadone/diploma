package depemployee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/depemployee"

	"github.com/lib/pq"
)

func (p PostgresEmployeeDepartment) UpdateEmployeeDepartment(
	ctx context.Context,
	sharedTx *sql.Tx,
	de *depemployee.DepartmentEmployee,
) error {
	if sharedTx == nil {
		return errors.New("transaction must be started before query")
	}

	query := `
        UPDATE employee_department
        SET
            position_id = $1,
            is_active = $2,
            updated_at = $3
        WHERE employee_id = $4
        AND department_id = $5
    `

	// Используем ExecContext вместо QueryRowContext для UPDATE
	result, err := sharedTx.ExecContext(
		ctx,
		query,
		de.PositionID,
		de.IsActive,
		de.UpdatedAt,
		de.EmployeeID,
		de.DepartmentID,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "employee_department_position_id_fkey":
				return fmt.Errorf("position %s does not exist", de.PositionID)
			}
		}
		return fmt.Errorf("failed to update employee department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee department link not found (employee: %s, department: %s)",
			de.EmployeeID, de.DepartmentID)
	}

	return nil
}
