package depposition

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/depposition"

	"github.com/google/uuid"
)

func (p PostgresDepPosition) GetDepartmentPositionById(
	ctx context.Context,
	sharedTx *sql.Tx,
	positionID uuid.UUID,
) (*depposition.DepPosition, error) {
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
        WHERE id = $1
    `

	var position depposition.DepPosition
	err := sharedTx.QueryRowContext(ctx, query, positionID).Scan(
		&position.Id,
		&position.DepartmentId,
		&position.Level,
		&position.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("department position with ID %s not found", positionID)
		}
		return nil, fmt.Errorf("failed to get department position: %w", err)
	}

	return &position, nil
}
