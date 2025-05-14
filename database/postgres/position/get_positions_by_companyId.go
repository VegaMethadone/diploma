package position

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/position"

	"github.com/google/uuid"
)

func (p PostgresPosition) GetPositionsByCompanyId(
	ctx context.Context,
	sharedTx *sql.Tx,
	companyId uuid.UUID,
) (*[]position.Position, error) {
	if sharedTx == nil {
		return nil, errors.New("start transaction before query")
	}

	query := `
        SELECT 
            id,
            company_id,
            lvl,
            name,
            is_active,
            created_at,
            updated_at
        FROM positions
        WHERE company_id = $1
        ORDER BY lvl ASC, created_at DESC
    `

	rows, err := sharedTx.QueryContext(ctx, query, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to query positions: %w", err)
	}
	defer rows.Close()

	var positions []position.Position
	for rows.Next() {
		var pos position.Position
		if err := rows.Scan(
			&pos.ID,
			&pos.CompanyID,
			&pos.Lvl,
			&pos.Name,
			&pos.IsActive,
			&pos.CreatedAt,
			&pos.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan position: %w", err)
		}
		positions = append(positions, pos)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &positions, nil
}
