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
