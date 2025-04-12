package department

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/department"

	"github.com/lib/pq"
)

func (p PostgresDepartment) CreateDepartment(
	ctx context.Context,
	sharedTx *sql.Tx,
	department *department.Department,
) error {
	if sharedTx == nil {
		return errors.New("transaction must be started before calling CreateDepartment")
	}

	query := `
        INSERT INTO departments (
            id,
            company_id,
            name,
            description,
            avatar_url,
            parent_id,
            created_at,
            updated_at,
            is_active
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err := sharedTx.ExecContext(
		ctx,
		query,
		department.ID,
		department.CompanyID,
		department.Name,
		department.Description,
		department.AvatarURL,
		department.ParentID,
		department.CreatedAt,
		department.UpdatedAt,
		department.IsActive,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "departments_company_id_name_key":
				return fmt.Errorf("department with name '%s' already exists in this company", department.Name)
			case "departments_parent_id_fkey":
				return fmt.Errorf("parent department %s does not exist", department.ParentID)
			}
		}
		return fmt.Errorf("failed to create department: %w", err)
	}

	return nil
}
