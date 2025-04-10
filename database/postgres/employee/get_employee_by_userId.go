package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/employee"

	"github.com/google/uuid"
)

func (p PostgresEmployee) GetEmployeeByUserId(
	ctx context.Context,
	sharedTx *sql.Tx,
	userId uuid.UUID,
) (*employee.Employee, error) {
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
        WHERE user_id = $1
        LIMIT 1
    `

	var empl employee.Employee
	err := sharedTx.QueryRowContext(ctx, query, userId).Scan(
		&empl.ID,
		&empl.UserID,
		&empl.CompanyID,
		&empl.PositionID,
		&empl.IsActive,
		&empl.IsOnline,
		&empl.LastActivityAt,
		&empl.CreatedAt,
		&empl.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee not found for user_id: %s", userId)
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	return &empl, nil
}
