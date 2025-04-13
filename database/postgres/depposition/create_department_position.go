package depposition

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/depposition"

	"github.com/lib/pq"
)

func (p PostgresDepPosition) CreateDepartmentPosition(
	ctx context.Context,
	sharedTx *sql.Tx,
	position *depposition.DepPosition,
) error {
	if sharedTx == nil {
		return errors.New("transaction must be started before query")
	}

	query := `
        INSERT INTO department_positions (
            id,
            department_id,
            level,
            name
        ) VALUES ($1, $2, $3, $4)
    `

	_, err := sharedTx.ExecContext(
		ctx,
		query,
		position.Id,
		position.DepartmentId,
		position.Level,
		position.Name,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "department_positions_pkey":
				return fmt.Errorf("position with ID %s already exists", position.Id)
			case "unique_position_name":
				return fmt.Errorf("position name '%s' already exists in this department", position.Name)
			case "department_positions_department_id_fkey":
				return fmt.Errorf("department %s does not exist", position.DepartmentId)
			}
		}
		return fmt.Errorf("failed to create department position: %w", err)
	}

	return nil
}
