package depemployee

import (
	"labyrinth/logic"

	"github.com/google/uuid"
)

const (
	userIDKey string = "id"
)

var bl *logic.BusinessLogic = logic.NewBusinessLogic()

type DepEmployeeHandlers struct{}

func NewDepEmployeeHandlers() DepEmployeeHandlers { return DepEmployeeHandlers{} }

type depemployeeData struct {
	EmployeeId uuid.UUID `json: "employee_id"`
	PositionId uuid.UUID `json: "position_id"`
}
