package depemployee

import (
	"fmt"
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

func (d *DepartmentEmployee) Print() {
	fmt.Println("======================")
	fmt.Printf("ID\t%s\nEmployeeID\t%s\nDepartmentID\t%s\nPositionID\t%s\n", d.ID, d.EmployeeID, d.DepartmentID, d.PositionID)
	fmt.Println("======================")
}
