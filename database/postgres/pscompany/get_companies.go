package postgres

import (
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/entity/company"

	"github.com/google/uuid"
)

func GetUserCompanies(userId uuid.UUID) ([]company.CompanyLogin, error) {
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	findQuery := `
        SELECT user_id, company_id
        FROM employee_company
        WHERE user_id = $1
    `
	rows, err := db.Query(findQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var userCompanies []company.CompanyLogin

	for rows.Next() {
		var uc company.CompanyLogin
		if err := rows.Scan(&uc.UserId, &uc.CompanyId); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err) // возможно тут вылезет ошибка, что нет строк, надо обработать
		}
		userCompanies = append(userCompanies, uc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return userCompanies, nil
}
