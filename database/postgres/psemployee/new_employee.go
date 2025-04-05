package postgres

import (
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"log"

	"github.com/google/uuid"
)

func NewEmployee(userId, companyId, positionId uuid.UUID) error {
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return err
	}
	defer db.Close()

	addQuery := `
		INSERT INTO employee_company
		(id, user_id, company_id, position)
		VALUES ($1, $2, $3, $4)
	`
	_, err = db.Exec(addQuery,
		uuid.New(),
		userId,
		companyId,
		positionId,
	)
	if err != nil {
		return fmt.Errorf("failed to insert new employee: %v", err)
	}

	log.Println("Employee registered successfully!")
	return nil
}
