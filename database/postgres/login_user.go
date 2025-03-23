package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (p *Postgres) LoginUser(login, password string) (*uuid.UUID, error) {
	db, err := sql.Open("postgres", p.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	findQuery := `
		SELECT id, password FROM users
		WHERE login = $1
	`

	var (
		userID         uuid.UUID
		hashedPassword string
	)

	err = db.QueryRow(findQuery, login).Scan(&userID, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &userID, nil
}
