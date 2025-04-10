package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/employee"

	"github.com/lib/pq"
)

func (p PostgresEmployee) UpdateEmployee(
	ctx context.Context,
	sharedTx *sql.Tx,
	empl *employee.Employee,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	query := `
        UPDATE employee_company
        SET
            user_id = $1,
            company_id = $2,
            position_id = $3,
            is_active = $4,
            is_online = $5,
            last_activity_at = $6,
            updated_at = $7
        WHERE id = $8
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		empl.UserID,
		empl.CompanyID,
		empl.PositionID,
		empl.IsActive,
		empl.IsOnline,
		empl.LastActivityAt,
		empl.UpdatedAt,
		empl.ID,
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
		return fmt.Errorf("failed to update employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found (id: %s)", empl.ID)
	}

	return nil
}
