package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/company"

	"github.com/lib/pq"
)

func (r PostgresCompany) UpdateCompany(
	ctx context.Context,
	sharedTx *sql.Tx,
	company *company.Company,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	query := `
        UPDATE companies
        SET
            owner_id = $1,
            name = $2,
            description = $3,
            logo_url = $4,
            industry = $5,
            employees = $6,
            is_verified = $7,
            is_active = $8,
            founded_date = $9,
            address = $10,
            phone = $11,
            email = $12,
            tax_number = $13,
            updated_at = NOW()
        WHERE id = $14
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		company.OwnerID,
		company.Name,
		company.Description,
		company.LogoURL,
		company.Industry,
		company.Employees,
		company.IsVerified,
		company.IsActive,
		company.FoundedDate,
		company.Address,
		company.Phone,
		company.Email,
		company.TaxNumber,
		company.ID,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "companies_name_key":
				return fmt.Errorf("company name already exists")
			case "companies_owner_id_fkey":
				return fmt.Errorf("owner does not exist")
			}
		}
		return fmt.Errorf("failed to update company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("company not found")
	}

	return nil
}
