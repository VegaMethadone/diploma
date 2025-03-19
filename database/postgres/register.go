package postgres

import (
	"database/sql"
	"fmt"
	"labyrinth/entiry/user"
	"log"

	_ "github.com/lib/pq"
)

func (p *Postgres) RegisterUser(user *user.User) error {
	db, err := sql.Open("postgres", p.conn)
	if err != nil {
		return err
	}
	defer db.Close()

	addQuery := `
		INSERT INTO users 
		(id, login, password, name, surname, bio, phone, telegram, mail)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = db.Exec(addQuery,
		user.Id,
		user.Login,
		user.Password,
		user.Name,
		user.Surname,
		user.Bio,
		user.Phone,
		user.Telegram,
		user.Mail,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	log.Println("User registered successfully!")
	return nil
}
