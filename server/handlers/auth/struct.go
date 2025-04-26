package auth

type userRegisterRequest struct {
	Mail         string `json: "mail"`
	HashPassword string `json: "password"`
	Phone        string `json: "phone"`
}
