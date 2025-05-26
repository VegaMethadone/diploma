package depemployee

import (
	"time"

	"github.com/google/uuid"
)

type DepartmentEmployee struct {
	ID           uuid.UUID `json:"id"`
	EmployeeID   uuid.UUID `json:"employee_id"`
	DepartmentID uuid.UUID `json:"department_id"`
	PositionID   uuid.UUID `json:"position_id"`
	CreatedAt    time.Time `json: "created_at"`
	UpdatedAt    time.Time `json: "updated_at"`
	IsActive     bool      `json:"is_active"`
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
