package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *PostgresCompany) DeactivateCompanyUsers(
	ctx context.Context,
	sharedTx *sql.Tx,
	companyID uuid.UUID,
) error {
	if sharedTx == nil {
		return errors.New("transaction is required")
	}

	query := `
        UPDATE user_companies
        SET isActive = false
        WHERE company_id = $1
        AND isActive = true
    `

	_, err := sharedTx.ExecContext(ctx, query, companyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("failed to deactivate company users: %w", err)
	}

	return nil
}
