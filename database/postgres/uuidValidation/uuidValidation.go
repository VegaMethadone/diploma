package uuidvalidation

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func checkUUID(ctx context.Context, db *sql.DB, id uuid.UUID) (bool, error) {
	const query = `SELECT 1 FROM used_uuids WHERE uuid_id = $1`

	var exists int
	err := db.QueryRowContext(ctx, query, id).Scan(&exists)

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
func reserveUUID(ctx context.Context, db *sql.DB, id uuid.UUID) error {
	const query = `INSERT INTO used_uuids (uuid_id) VALUES ($1)`

	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to reserve UUID: %w", err)
	}

	return nil
}

func CheckAndReserverUUID(ctx context.Context, db *sql.DB, newUUIDfunc func() uuid.UUID) (uuid.UUID, error) {
	var exists bool = true
	var newUUID = uuid.Nil
	for !exists {
		generatedUUID := newUUIDfunc()
		ok, err := checkUUID(ctx, db, generatedUUID)
		if err != nil {
			return uuid.Nil, err // верни ошибку,  что  зафелился CheckUUID
		}
		if ok {
			newUUID = generatedUUID
			exists = false
		}
	}
	err := reserveUUID(ctx, db, newUUID)
	if err != nil {
		return uuid.Nil, err
	}

	return newUUID, nil
}

// func CheckAndReserveUUID(ctx context.Context, db *sql.DB, generateUUID func() uuid.UUID) (uuid.UUID, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return uuid.Nil, fmt.Errorf("таймаут при поиске свободного UUID")
// 		default:
// 			// Генерируем новый UUID
// 			newUUID := generateUUID()

// 			// Проверяем доступность
// 			free, err := CheckUUID(ctx, db, newUUID)
// 			if err != nil {
// 				continue // Пропускаем ошибки проверки
// 			}

// 			if free {
// 				// Пытаемся зарезервировать
// 				if err := ReserveUUID(ctx, db, newUUID); err == nil {
// 					return newUUID, nil
// 				}
// 			}
// 		}
// 	}
// }
