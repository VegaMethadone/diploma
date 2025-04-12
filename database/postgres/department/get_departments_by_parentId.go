package department

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/department"

	"github.com/google/uuid"
)

func (p PostgresDepartment) GetDepartmentsByParentId(
	ctx context.Context,
	sharedTx *sql.Tx,
	parentId uuid.UUID,
) ([]*department.Department, error) {
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
        WHERE parent_id = $1
        AND is_active = true
        ORDER BY name ASC
    `

	rows, err := sharedTx.QueryContext(ctx, query, parentId)
	if err != nil {
		return nil, fmt.Errorf("failed to query departments: %w", err)
	}
	defer rows.Close()

	var departments []*department.Department
	for rows.Next() {
		var dept department.Department
		if err := rows.Scan(
			&dept.ID,
			&dept.CompanyID,
			&dept.Name,
			&dept.Description,
			&dept.AvatarURL,
			&dept.ParentID,
			&dept.CreatedAt,
			&dept.UpdatedAt,
			&dept.IsActive,
		); err != nil {
			return nil, fmt.Errorf("failed to scan department: %w", err)
		}
		departments = append(departments, &dept)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return departments, nil
}
