package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *PostgresCompany) DeleteCompany(
	ctx context.Context,
	sharedTx *sql.Tx,
	id uuid.UUID,
) error {
	if sharedTx == nil {
		return errors.New("transaction is required")
	}

	query := `
        UPDATE companies 
        SET is_active = false,
            updated_at = NOW()
        WHERE id = $1
        AND is_active = true
    `

	result, err := sharedTx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("company not found or already deleted")
	}

	return nil
}
