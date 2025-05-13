package uuidvalidation

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type DBUuidValidation struct{}

func NewDBUuidValidation() DBUuidValidation { return DBUuidValidation{} }

func checkUUID(ctx context.Context, sharedTx *sql.Tx, id uuid.UUID) (bool, error) {
	const query = `SELECT 1 FROM used_uuids WHERE uuid_id = $1`

	var exists int
	err := sharedTx.QueryRowContext(ctx, query, id).Scan(&exists)

	switch {
	case err == sql.ErrNoRows:
		return true, nil // UUID не найден - свободен
	case err != nil:
		return false, fmt.Errorf("failed to check UUID: %w", err)
	default:
		return false, nil // UUID найден - занят
	}
}

// ReserveUUID добавляет UUID в таблицу used_uuids
func reserveUUID(ctx context.Context, sharedTx *sql.Tx, id uuid.UUID) error {
	const query = `INSERT INTO used_uuids (uuid_id) VALUES ($1)`

	_, err := sharedTx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to reserve UUID: %w", err)
	}

	return nil
}

func (u DBUuidValidation) CheckAndReserveUUID(ctx context.Context, sharedTx *sql.Tx) (uuid.UUID, error) {
	var exists bool = true
	var newUUID = uuid.Nil
	for exists {
		generatedUUID := uuid.New()
		ok, err := checkUUID(ctx, sharedTx, generatedUUID)
		if err != nil {
			return uuid.Nil, err
		}
		if ok {
			newUUID = generatedUUID
			exists = false
		}
	}
	err := reserveUUID(ctx, sharedTx, newUUID)
	if err != nil {
		return uuid.Nil, err
	}

	return newUUID, nil
}
