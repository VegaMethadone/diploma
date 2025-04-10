package position

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/position"
	"time"

	"github.com/lib/pq"
)

func (p PostgresPosition) UpdatePosition(
	ctx context.Context,
	sharedTx *sql.Tx,
	position *position.Position,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	query := `
        UPDATE positions
        SET
            company_id = $1,
            lvl = $2,
            name = $3,
            is_active = $4,
            updated_at = $5
        WHERE id = $6
    `

	position.UpdatedAt = time.Now()

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		position.CompanyID,
		position.Lvl,
		position.Name,
		position.IsActive,
		position.UpdatedAt,
		position.ID,
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
		return fmt.Errorf("failed to update position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("position not found (id: %s)", position.ID)
	}

	return nil
}
