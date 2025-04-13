package depemployee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/depemployee"

	"github.com/google/uuid"
)

func (p PostgresEmployeeDepartment) GetEmployeesDepartmentByDepartmentId(
	ctx context.Context,
	sharedTx *sql.Tx,
	departmentID uuid.UUID,
) ([]*depemployee.DepartmentEmployee, error) {
	if sharedTx == nil {
		return nil, errors.New("transaction must be started before query")
	}

	query := `
        SELECT 
            id,
            employee_id,
            department_id,
            position_id,
            is_active,
            created_at,
            updated_at
        FROM employee_department
        WHERE department_id = $1
        ORDER BY created_at DESC
    `

	rows, err := sharedTx.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query department employees: %w", err)
	}
	defer rows.Close()

	var employees []*depemployee.DepartmentEmployee
	for rows.Next() {
		var de depemployee.DepartmentEmployee
		if err := rows.Scan(
			&de.ID,
			&de.EmployeeID,
			&de.DepartmentID,
			&de.PositionID,
			&de.IsActive,
			&de.CreatedAt,
			&de.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan employee department: %w", err)
		}
		employees = append(employees, &de)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return employees, nil
}
