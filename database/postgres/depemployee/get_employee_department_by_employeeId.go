package depemployee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/depemployee"

	"github.com/google/uuid"
)

func (p PostgresEmployeeDepartment) GetEmployeeDepartmentByEmployeeId(
	ctx context.Context,
	sharedTx *sql.Tx,
	employeeID uuid.UUID,
	departmentID uuid.UUID,
) (*depemployee.DepartmentEmployee, error) {
	if sharedTx == nil {
		return nil, errors.New("transaction must be started before query")
	}

	query := `
        SELECT 
            id,
            employee_id,
            department_id,
            position_id,
            created_at,
            updated_at,
            is_active
        FROM employee_department
        WHERE employee_id = $1
        AND department_id = $2
        LIMIT 1
    `

	var de depemployee.DepartmentEmployee
	err := sharedTx.QueryRowContext(ctx, query, employeeID, departmentID).Scan(
		&de.ID,
		&de.EmployeeID,
		&de.DepartmentID,
		&de.PositionID,
		&de.CreatedAt,
		&de.UpdatedAt,
		&de.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee department link not found (employee: %s, department: %s)",
				employeeID, departmentID)
		}
		return nil, fmt.Errorf("failed to get employee department: %w", err)
	}

	return &de, nil
}
