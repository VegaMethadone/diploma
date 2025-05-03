package company

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID          uuid.UUID `json:"id"`           // Уникальный идентификатор
	OwnerID     uuid.UUID `json:"owner_id"`     // ID владельца (лучше явно указать _id)
	Name        string    `json:"name"`         // Название компании (обязательное)
	Description string    `json:"description"`  // Описание
	LogoURL     string    `json:"logo_url"`     // Ссылка на логотип
	Industry    string    `json:"industry"`     // Отрасль
	Employees   int       `json:"employees"`    // Количество сотрудников
	IsVerified  bool      `json:"is_verified"`  // Подтверждена ли компания
	IsActive    bool      `json:"is_active"`    // Активна ли компания
	CreatedAt   time.Time `json:"created_at"`   // Дата создания (авто)
	UpdatedAt   time.Time `json:"updated_at"`   // Дата обновления (авто)
	FoundedDate time.Time `json:"founded_date"` // Дата основания
	Address     string    `json:"address"`      // Адрес
	Phone       string    `json:"phone"`        // Телефон
	Email       string    `json:"email"`        // Email
	TaxNumber   string    `json:"tax_number"`   // Добавлено: налоговый номер
}

func NewCompany(ownerId uuid.UUID, name, description, address, phone, email string) *Company {
	return &Company{
		ID:          uuid.New(),
		OwnerID:     ownerId,
		Name:        name,
		Description: description,
		LogoURL:     "",
		Industry:    "",
		Employees:   1,
		IsVerified:  false,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		FoundedDate: time.Now(),
		Address:     address,
		Phone:       phone,
		Email:       email,
		TaxNumber:   "",
	}
}
