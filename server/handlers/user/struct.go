package user

import (
	"labyrinth/logic"
	"labyrinth/models/user"
	"time"

	"github.com/google/uuid"
)

var bl logic.BusinessLogic = *logic.NewBusinessLogic()

type UserHandlers struct{}

func NewUserHandlers() UserHandlers { return UserHandlers{} }

type userData struct {
	Login            string    `json:"login"`             // Уникальный логин
	PasswordHash     string    `json:"password"`          // Хеш пароля // нельзя передавать хеш в json
	Email            string    `json:"email"`             // Уникальный email
	EmailVerified    bool      `json:"email_verified"`    // Подтвержден ли email
	Phone            string    `json:"phone"`             // Телефон с валидацией
	PhoneVerified    bool      `json:"phone_verified"`    // Подтвержден ли телефон
	FirstName        string    `json:"first_name"`        // Имя
	LastName         string    `json:"last_name"`         // Фамилия
	Bio              string    `json:"bio"`               // Биография
	TelegramUsername string    `json:"telegram_username"` // Telegram
	AvatarURL        string    `json:"avatar_url"`        // Ссылка на аватар
	CreatedAt        time.Time `json:"created_at"`        // Дата создания
}

func newUserData(usr *user.User) *userData {
	return &userData{
		Login:            usr.Login,
		PasswordHash:     usr.PasswordHash,
		Email:            usr.Email,
		EmailVerified:    usr.EmailVerified,
		Phone:            usr.Phone,
		PhoneVerified:    usr.PhoneVerified,
		FirstName:        usr.FirstName,
		LastName:         usr.LastName,
		Bio:              usr.Bio,
		TelegramUsername: usr.TelegramUsername,
		AvatarURL:        usr.AvatarURL,
		CreatedAt:        usr.CreatedAt,
	}
}

func NewUser(usr *userData, id uuid.UUID) *user.User {
	return &user.User{
		ID:               id,
		Login:            usr.Login,
		PasswordHash:     usr.PasswordHash,
		Email:            usr.Email,
		EmailVerified:    usr.EmailVerified,
		Phone:            usr.Phone,
		PhoneVerified:    usr.PhoneVerified,
		FirstName:        usr.FirstName,
		LastName:         usr.LastName,
		Bio:              usr.Bio,
		TelegramUsername: usr.TelegramUsername,
		AvatarURL:        usr.AvatarURL,
		CreatedAt:        usr.CreatedAt,
		UpdatedAt:        time.Now(),
		LastLoginAt:      time.Now(),
		IsActive:         true,
		IsStaff:          false,
	}
}
