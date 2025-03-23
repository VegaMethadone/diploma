package user

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Login    string
	Password string
	Name     string
	Surname  string
	Bio      string
	Phone    string
	Telegram string
	Mail     string
}

func NewUser(login, passwrod string) *User {
	return &User{
		Id:       uuid.New(),
		Login:    login,
		Password: passwrod,
		Name:     "",
		Surname:  "",
		Bio:      "",
		Phone:    "",
		Telegram: "",
		Mail:     login,
	}
}
