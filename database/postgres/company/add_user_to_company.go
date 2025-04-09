package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (r *PostgresCompany) AddUserToCompany(
	ctx context.Context,
	sharedTx *sql.Tx,
	userID uuid.UUID,
	companyID uuid.UUID,
) error {
	if sharedTx == nil {
		return errors.New("transaction is required")
	}

	query := `
		INSERT INTO user_companies (
            user_id,
            company_id,
            isActive
        ) VALUES ($1, $2, $3)
	`
	_, err := sharedTx.ExecContext(
		ctx,
		query,
		userID,
		companyID,
		true,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "user_companies_user_id_fkey":
				return fmt.Errorf("user does not exist")
			case "user_companies_company_id_fkey":
				return fmt.Errorf("company does not exist")
			}
		}
		return fmt.Errorf("failed to add user to company: %w", err)
	}

	return nil
}
