package depemployee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/depemployee"

	"github.com/lib/pq"
)

func (p PostgresEmployeeDepartment) CreateEmployeeDepartment(
	ctx context.Context,
	sharedTx *sql.Tx,
	de *depemployee.DepartmentEmployee,
) error {
	if sharedTx == nil {
		return errors.New("transaction must be started before calling CreateEmployeeDepartment")
	}

	query := `
        INSERT INTO employee_department (
            id,
            employee_id,
            department_id,
            position_id,
            is_active,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := sharedTx.ExecContext(
		ctx,
		query,
		de.ID,
		de.EmployeeID,
		de.DepartmentID,
		de.PositionID,
		de.IsActive,
		de.CreatedAt,
		de.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "employee_department_employee_id_department_id_key":
				return fmt.Errorf("employee %s already exists in department %s", de.EmployeeID, de.DepartmentID)
			case "employee_department_employee_id_fkey":
				return fmt.Errorf("employee %s does not exist", de.EmployeeID)
			case "employee_department_department_id_fkey":
				return fmt.Errorf("department %s does not exist", de.DepartmentID)
			case "employee_department_position_id_fkey":
				return fmt.Errorf("position %s does not exist", de.PositionID)
			}
		}
		return fmt.Errorf("failed to create employee department: %w", err)
	}

	return nil
}
