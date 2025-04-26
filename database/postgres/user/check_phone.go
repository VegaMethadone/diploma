package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (p PostgresUser) CheckPhone(ctx context.Context, sharedTx *sql.Tx, phone string) (bool, error) {
	if sharedTx == nil {
		return false, errors.New("start transaction before query")
	}

	if phone == "" {
		return false, errors.New("phone number cannot be empty")
	}

	query := `
        SELECT EXISTS (
            SELECT 1 FROM users
            WHERE phone = $1 AND phone_verified = true
        )
    `

	var exists bool
	err := sharedTx.QueryRowContext(ctx, query, phone).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check phone: %w", err)
	}

	return exists, nil
}
