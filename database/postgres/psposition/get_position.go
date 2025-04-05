package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/entity/position"

	"github.com/google/uuid"
)

func GetPosition(positionId uuid.UUID) (*position.Position, error) {
	// Проверка на нулевой UUID
	if positionId == uuid.Nil {
		return nil, fmt.Errorf("invalid position ID")
	}

	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Явное перечисление полей вместо SELECT *
	getQuery := `
        SELECT id, company_id, lvl, name 
        FROM positions
        WHERE id = $1
    `

	var pos position.Position
	err = db.QueryRow(getQuery, positionId).Scan(
		&pos.Id,
		&pos.CompanyId,
		&pos.Lvl,
		&pos.Name,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		// Позиция не найдена
		return nil, nil
	case err != nil:
		// Другие ошибки при выполнении запроса
		return nil, fmt.Errorf("failed to get position: %w", err)
	default:
		// Успешное выполнение
		return &pos, nil
	}
}
