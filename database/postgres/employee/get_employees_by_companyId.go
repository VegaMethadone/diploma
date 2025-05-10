package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/employee"

	"github.com/google/uuid"
)

func (p PostgresEmployee) GetEmployeesByCompanyId(
	ctx context.Context,
	sharedTx *sql.Tx,
	companyId uuid.UUID,
) (*[]employee.Employee, error) {
	if sharedTx == nil {
		return nil, errors.New("start transaction before query")
	}

	query := `
        SELECT 
            id,
            user_id,
            company_id,
            position_id,
            is_active,
            is_online,
            last_activity_at,
            created_at,
            updated_at
        FROM employee_company
        WHERE company_id = $1 and is_active = true
        ORDER BY created_at DESC
    `

	rows, err := sharedTx.QueryContext(ctx, query, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to query employees: %w", err)
	}
	defer rows.Close()

	var employees []employee.Employee
	for rows.Next() {
		var empl employee.Employee
		if err := rows.Scan(
			&empl.ID,
			&empl.UserID,
			&empl.CompanyID,
			&empl.PositionID,
			&empl.IsActive,
			&empl.IsOnline,
			&empl.LastActivityAt,
			&empl.CreatedAt,
			&empl.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, empl)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &employees, nil
}
