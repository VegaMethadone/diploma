package department

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/department"

	"github.com/google/uuid"
)

func (p PostgresDepartment) GetDepartmentById(
	ctx context.Context,
	sharedTx *sql.Tx,
	departmentId uuid.UUID,
) (*department.Department, error) {
	if sharedTx == nil {
		return nil, errors.New("start transaction before query")
	}

	query := `
        SELECT 
            id,
            company_id,
            name,
            description,
            avatar_url,
            parent_id,
            created_at,
            updated_at,
            is_active
        FROM departments
        WHERE id = $1
        LIMIT 1
    `

	var dept department.Department
	err := sharedTx.QueryRowContext(ctx, query, departmentId).Scan(
		&dept.ID,
		&dept.CompanyID,
		&dept.Name,
		&dept.Description,
		&dept.AvatarURL,
		&dept.ParentID,
		&dept.CreatedAt,
		&dept.UpdatedAt,
		&dept.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("department not found (id: %s)", departmentId)
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return &dept, nil
}
