package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"

	"github.com/google/uuid"
)

func CheckEmployee(userId, employeeId, companyId uuid.UUID) (bool, error) {
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return false, fmt.Errorf("database connection error: %w", err)
	}
	defer db.Close()

	checkQuery := `
        SELECT user_id
        FROM employee_company
        WHERE id = $1 AND company_id = $2
    `

	var gotUserId uuid.UUID
	err = db.QueryRow(checkQuery, employeeId, companyId).Scan(&gotUserId)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		// Сотрудник не найден в указанной компании
		return false, nil
	case err != nil:
		// Произошла другая ошибка при выполнении запроса
		return false, fmt.Errorf("database query error: %w", err)
	default:
		// Запрос выполнен успешно, сравниваем user_id
		return userId == gotUserId, nil
	}
}
