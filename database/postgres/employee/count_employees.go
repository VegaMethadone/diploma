package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (p PostgresEmployee) CountEmployees(
	ctx context.Context,
	sharedTx *sql.Tx,
	companyId uuid.UUID,
) (int, error) {
	if sharedTx == nil {
		return 0, errors.New("start transaction before query")
	}

	query := `
        SELECT COUNT(*)
        FROM employee_company
        WHERE company_id = $1
        AND is_active = true
    `

	var count int
	err := sharedTx.QueryRowContext(ctx, query, companyId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count employees: %w", err)
	}

	return count, nil
}
