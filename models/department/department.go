package department

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	ID          uuid.UUID `json:"id"`
	CompanyID   uuid.UUID `json:"company_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatar_url"` // Добавлено из таблицы
	ParentID    uuid.UUID `json:"parent"`     // Соответствует полю parent_id в таблице
	CreatedAt   time.Time `json:"created_at"` // Добавлено из таблицы
	UpdatedAt   time.Time `json:"updated_at"` // Добавлено из таблицы
	IsActive    bool      `json:"is_active"`  // Добавлено из таблицы
}

func NewDepartment(
	generatedId,
	companyId,
	parentId uuid.UUID,
	name,
	description string,
) *Department {
	return &Department{
		ID:          generatedId,
		CompanyID:   companyId,
		Name:        name,
		Description: description,
		AvatarURL:   "",
		ParentID:    parentId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
	}
}
