package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/company"

	"github.com/lib/pq"
)

func (r PostgresCompany) CreateCompany(
	ctx context.Context,
	sharedTx *sql.Tx,
	company *company.Company,
) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}

	query := `
        INSERT INTO companies (
            id,
            owner_id,
            name,
            description,
            logo_url,
            industry,
            employees,
            is_verified,
            is_active,
            founded_date,
            address,
            phone,
            email,
            tax_number
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		company.ID,
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
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "companies_pkey":
				return fmt.Errorf("company with this ID already exists")
			case "companies_name_key":
				return fmt.Errorf("company with this name already exists")
			}
		}
		return fmt.Errorf("failed to create company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found (id: %s)", company.ID)
	}

	return nil
}
