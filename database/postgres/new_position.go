package postgres

import (
	"database/sql"
	"fmt"
	"labyrinth/entity/position"
)

func (p *Postgres) NewPosition(pos *position.Position) error {
	db, err := sql.Open("postgres", p.conn)
	if err != nil {
		return err
	}
	defer db.Close()

	addQuery := `
        INSERT INTO positions (id, company_id, lvl, name)
        VALUES ($1, $2, $3, $4)
    `

	_, err = db.Exec(addQuery, pos.Id, pos.CompanyId, pos.Lvl, pos.Name)
	if err != nil {
		return fmt.Errorf("failed to create position: %w", err)
	}

	return nil

}
