package postgres

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func (p *Postgres) CheckUser(userId uuid.UUID) (bool, error) {
	db, err := sql.Open("postgres", p.conn)
	if err != nil {
		return false, err
	}
	defer db.Close()

	findQuery := `
		SELECT id FROM users
		WHERE id = $1
	`

	var foundUderId uuid.UUID
	err = db.QueryRow(findQuery, userId).Scan(&foundUderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return (userId == foundUderId), nil
}
