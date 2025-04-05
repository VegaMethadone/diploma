package postgres

import (
	"database/sql"
	"errors"
	"labyrinth/database/postgres"

	"github.com/google/uuid"
)

func CheckCompany(companyId uuid.UUID) (bool, error) {
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return false, err
	}
	defer db.Close()

	findQuery := `
		SELECT id FROM companies
		WHERE id == $1
	`

	var foundCompanyId uuid.UUID
	err = db.QueryRow(findQuery, companyId).Scan(&foundCompanyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return (companyId == foundCompanyId), nil
}
