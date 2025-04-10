package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/employee"

	"github.com/lib/pq"
)

func (p PostgresEmployee) CreateEmployee(
	ctx context.Context,
	sharedTx *sql.Tx,
	empl *employee.Employee,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	query := `
        INSERT INTO employee_company (
            id,
            user_id,
            company_id,
            position_id,
            is_active,
            is_online,
            last_activity_at,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err := sharedTx.ExecContext(
		ctx,
		query,
		empl.ID,
		empl.UserID,
		empl.CompanyID,
		empl.PositionID,
		empl.IsActive,
		empl.IsOnline,
		empl.LastActivityAt,
		empl.CreatedAt,
		empl.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "employee_company_user_id_company_id_key":
				return fmt.Errorf("user %s already exists in company %s", empl.UserID, empl.CompanyID)
			case "employee_company_position_id_fkey":
				return fmt.Errorf("position %s does not exist", empl.PositionID)
			}
		}
		return fmt.Errorf("failed to create employee: %w", err)
	}

	return nil
}
