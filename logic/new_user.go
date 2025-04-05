package logic

import (
	psuser "labyrinth/database/postgres/psuser"
	"labyrinth/entity/user"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func newUser(mail, passwrod string) (*user.User, error) {
	hashedPassword, err := hashPassword(passwrod)
	if err != nil {
		return nil, err
	}

	return user.NewUser(mail, hashedPassword), nil
}

func NewUser(mail, password string) error {
	user, err := newUser(mail, password)
	if err != nil {
		return err
	}

	return psuser.RegisterUser(user)
}
