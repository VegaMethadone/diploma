package depposition

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/depposition"

	"github.com/google/uuid"
)

func (p PostgresDepPosition) GetDepartmentPositionsByDepartmentId(
	ctx context.Context,
	sharedTx *sql.Tx,
	departmentID uuid.UUID,
) (*[]depposition.DepPosition, error) {
	if sharedTx == nil {
		return nil, errors.New("transaction must be started before query")
	}

	query := `
        SELECT 
            id,
            department_id,
            level,
            name
        FROM department_positions
        WHERE department_id = $1
        ORDER BY level, name
    `

	rows, err := sharedTx.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query department positions: %w", err)
	}
	defer rows.Close()

	var positions []depposition.DepPosition
	for rows.Next() {
		var position depposition.DepPosition
		if err := rows.Scan(
			&position.Id,
			&position.DepartmentId,
			&position.Level,
			&position.Name,
		); err != nil {
			return nil, fmt.Errorf("failed to scan department position: %w", err)
		}
		positions = append(positions, position)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating department positions: %w", err)
	}

	return &positions, nil
}
