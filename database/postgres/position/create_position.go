package position

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/position"

	"github.com/lib/pq"
)

func (p PostgresPosition) CreatePosition(
	ctx context.Context,
	sharedTx *sql.Tx,
	position *position.Position,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	query := `
        INSERT INTO positions (
            id,
            company_id,
            lvl,
            name,
            is_active,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		position.ID,
		position.CompanyID,
		position.Lvl,
		position.Name,
		position.IsActive,
		position.CreatedAt,
		position.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "positions_company_id_name_key":
				return errors.New("position name must be unique within company")
			case "positions_company_id_lvl_key":
				return fmt.Errorf("position level %d already exists in company", position.Lvl)
			}
		}
		return fmt.Errorf("failed to create position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected (position id: %s)", position.ID)
	}

	return nil
}
