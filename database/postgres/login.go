package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (p *Postgres) LoginUser(login, password string) error {
	db, err := sql.Open("postgres", p.conn)
	if err != nil {
		return err
	}
	defer db.Close()

	findQuery := `
		SELECT password FROM users
		WHERE login = $1
	`
	var hashedPassword string
	err = db.QueryRow(findQuery, login).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with login '%s' not found", login)
		}
		return fmt.Errorf("failed to query user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}
