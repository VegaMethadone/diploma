package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (p PostgresEmployee) ExistsEmployee(
	ctx context.Context,
	sharedTx *sql.Tx,
	employeeId uuid.UUID,
) (bool, error) {
	if sharedTx == nil {
		return false, errors.New("start transaction before query")
	}

	query := `
        SELECT EXISTS(
            SELECT 1 
            FROM employee_company 
            WHERE id = $1
        )
    `

	var exists bool
	err := sharedTx.QueryRowContext(ctx, query, employeeId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check employee existence: %w", err)
	}

	return exists, nil
}
