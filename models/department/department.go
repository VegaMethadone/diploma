package department

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	ID          uuid.UUID `json:"id"`
	CompanyID   uuid.UUID `json:"companyId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatarUrl"` // Добавлено из таблицы
	ParentID    uuid.UUID `json:"parent"`    // Соответствует полю parent_id в таблице
	CreatedAt   time.Time `json:"createdAt"` // Добавлено из таблицы
	UpdatedAt   time.Time `json:"updatedAt"` // Добавлено из таблицы
	IsActive    bool      `json:"isActive"`  // Добавлено из таблицы
}
