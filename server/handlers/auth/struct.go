package auth

import "labyrinth/logic"

var bl *logic.BusinessLogic = logic.NewBusinessLogic()

type userRegisterRequest struct {
	Mail         string `json: "mail"`
	HashPassword string `json: "password"`
	Phone        string `json: "phone"`
}

type userLoginRequest struct {
	Mail         string `json: "mail"`
	HashPassword string `json: "password"`
}
