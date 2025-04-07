package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID `json:"id"`
	Login            string    `json:"login"`             // Уникальный логин
	PasswordHash     string    `json:"-"`                 // Хеш пароля (исключен из JSON)
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

func NewUser(login, email, passwordHash string) *User {
	return &User{
		ID:           uuid.New(), // Генерация нового UUID
		Login:        strings.TrimSpace(login),
		PasswordHash: passwordHash, // Должен приходить уже хешированным
		Email:        strings.TrimSpace(email),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		IsActive:     true, // Новые пользователи активны по умолчанию
	}
}

// Validate проверяет обязательные поля
func (u *User) Validate() error {
	if u.Login == "" {
		return errors.New("login is required")
	}
	if u.PasswordHash == "" {
		return errors.New("password hash is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	return nil
}
