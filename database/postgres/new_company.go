package postgres

import (
	"database/sql"
	"fmt"
	"labyrinth/entity/company"
	"log"

	"github.com/google/uuid"
)

func (p *Postgres) RegisterCompany(company_ *company.Company) (*uuid.UUID, error) {
	db, err := sql.Open("postgres", p.conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	addQuery := `
		INSERT INTO companies
		(id, owner_id, text, description)
		VALUES ($1, $2, $3, $4)
	`

	_, err = db.Exec(addQuery,
		company_.Id,
		company_.Owner,
		company_.Text,
		company_.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert new company: %v", err)
	}

	log.Println("Company registered successfully!")
	return &company_.Id, nil
}
