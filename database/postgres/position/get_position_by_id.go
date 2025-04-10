package position

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/position"

	"github.com/google/uuid"
)

func (p PostgresPosition) GetPositionById(
	ctx context.Context,
	sharedTx *sql.Tx,
	positionId uuid.UUID,
) (*position.Position, error) {
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
        WHERE id = $1
    `

	var pos position.Position
	err := sharedTx.QueryRowContext(ctx, query, positionId).Scan(
		&pos.ID,
		&pos.CompanyID,
		&pos.Lvl,
		&pos.Name,
		&pos.IsActive,
		&pos.CreatedAt,
		&pos.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("position not found (id: %s)", positionId)
		}
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	return &pos, nil
}
