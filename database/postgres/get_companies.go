package postgres

import (
	"database/sql"
	"labyrinth/entity/company"

	"github.com/google/uuid"
)

func (p *Postgres) GetUserCompanies(uderUUID_ uuid.UUID) (*[]company.Company, error) {
	db, err := sql.Open("postgres", p.conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	findAsOwnerQuery := `
		SELECT id, owner_id, text
		FROM companies where owner_id = $1 
	`

	/*
		поменять  логику, что  я  буду  добавлять овнера сразу как имплоя
	*/
}
