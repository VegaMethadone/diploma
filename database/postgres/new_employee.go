package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func (p *Postgres) NewEmployee(userId, companyId uuid.UUID) error {
	db, err := sql.Open("postgres", p.conn)
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
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to insert new employee: %v", err)
	}

	log.Println("Employee registered successfully!")
	return nil
} // добавить потом nil в роли, либо добавить роли
