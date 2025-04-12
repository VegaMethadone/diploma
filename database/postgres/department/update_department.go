package department

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/department"
	"time"

	"github.com/lib/pq"
)

func (p PostgresDepartment) UpdateDepartment(
	ctx context.Context,
	sharedTx *sql.Tx,
	department *department.Department,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	// Обновляем метку времени
	department.UpdatedAt = time.Now()

	query := `
        UPDATE departments
        SET
            company_id = $1,
            name = $2,
            description = $3,
            avatar_url = $4,
            parent_id = $5,
            updated_at = $6,
            is_active = $7
        WHERE id = $8
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		department.CompanyID,
		department.Name,
		department.Description,
		department.AvatarURL,
		department.ParentID,
		department.UpdatedAt,
		department.IsActive,
		department.ID,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "departments_company_id_name_key":
				return fmt.Errorf("department name '%s' already exists in company", department.Name)
			case "departments_parent_id_fkey":
				return fmt.Errorf("parent department %s does not exist", department.ParentID)
			}
		}
		return fmt.Errorf("failed to update department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found (id: %s)", department.ID)
	}

	return nil
}
