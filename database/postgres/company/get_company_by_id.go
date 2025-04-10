package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/company"

	"github.com/google/uuid"
)

func (r PostgresCompany) GetCompanyByID(
	ctx context.Context,
	sharedTx *sql.Tx,
	id uuid.UUID,
) (*company.Company, error) {
	if sharedTx == nil {
		return nil, errors.New("start transaction before query")
	}

	query := `
        SELECT 
            id,
            owner_id,
            name,
            description,
            logo_url,
            industry,
            employees,
            is_verified,
            is_active,
            created_at,
            updated_at,
            founded_date,
            address,
            phone,
            email,
            tax_number
        FROM companies
        WHERE id = $1
        LIMIT 1
    `

	var c company.Company

	err := sharedTx.QueryRowContext(ctx, query, id).Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("company not found")
		}
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	return &c, nil
}
