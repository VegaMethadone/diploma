package depemployee

import (
	"time"

	"github.com/google/uuid"
)

type DepartmentEmployee struct {
	ID           uuid.UUID `json:"id"`
	EmployeeID   uuid.UUID `json:"employeeId"`
	DepartmentID uuid.UUID `json:"departmentId"`
	PositionID   uuid.UUID `json:"positionId"` // Соответствует полю position_id из таблицы
	CreatedAt    time.Time `json: "createdAt"`
	UpdatedAt    time.Time `json: "updatedAt"`
	IsActive     bool      `json:"isActive"` // Добавлено из таблицы
}

func NewDepartmentEmployee(
	generatedId,
	employeeId,
	departmentId,
	positionId uuid.UUID,
) *DepartmentEmployee {
	return &DepartmentEmployee{
		ID:           generatedId,
		EmployeeID:   employeeId,
		DepartmentID: departmentId,
		PositionID:   positionId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}
}
