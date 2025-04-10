package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/company"

	"github.com/google/uuid"
)

func (r PostgresCompany) GetCompaniesByUser(
	ctx context.Context,
	sharedTx *sql.Tx,
	userID uuid.UUID,
) ([]*company.Company, error) {
	if sharedTx == nil {
		return nil, errors.New("start transaction before query")
	}

	query := `
        SELECT 
            c.id,
            c.owner_id,
            c.name,
            c.description,
            c.logo_url,
            c.industry,
            c.employees,
            c.is_verified,
            c.is_active,
            c.created_at,
            c.updated_at,
            c.founded_date,
            c.address,
            c.phone,
            c.email,
            c.tax_number
        FROM companies c
        JOIN user_companies uc ON c.id = uc.company_id
        WHERE uc.user_id = $1 
        AND uc.isActive = true
        AND c.is_active = true
    `

	rows, err := sharedTx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user companies: %w", err)
	}
	defer rows.Close()

	var companies []*company.Company

	for rows.Next() {
		var c company.Company

		err := rows.Scan(
			&c.ID,
			&c.OwnerID,
			&c.Name,
			&c.Description,
			&c.LogoURL,
			&c.Industry,
			&c.Employees,
			&c.IsVerified,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.FoundedDate,
			&c.Address,
			&c.Phone,
			&c.Email,
			&c.TaxNumber,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}

		companies = append(companies, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return companies, nil
}
