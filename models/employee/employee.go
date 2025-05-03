package employee

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	CompanyID      uuid.UUID `json:"company_id"`
	PositionID     uuid.UUID `json:"position_id"`
	IsActive       bool      `json:"is_active"`
	IsOnline       bool      `json:"is_online"`
	LastActivityAt time.Time `json:"last_activity_at,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewEmployee(
	generatedId,
	userId,
	companyId,
	positionId uuid.UUID,
) *Employee {
	return &Employee{
		ID:             generatedId,
		UserID:         userId,
		CompanyID:      companyId,
		PositionID:     positionId,
		IsActive:       true,
		IsOnline:       false,
		LastActivityAt: time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
