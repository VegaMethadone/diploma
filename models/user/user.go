package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID `json:"id"`
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
	UpdatedAt        time.Time `json:"updated_at"`        // Дата обновления
	LastLoginAt      time.Time `json:"last_login_at"`     // Последний вход
	IsActive         bool      `json:"is_active"`         // Активен ли аккаунт
	IsStaff          bool      `json:"is_staff"`          // Персонал/админ
}
