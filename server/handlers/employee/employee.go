package employee

import (
	"labyrinth/logic"

	"github.com/google/uuid"
)

const (
	userIDKey string = "id"
)

var bl *logic.BusinessLogic = logic.NewBusinessLogic()

type EmployeeHandlers struct{}

func NewEmployeeHandlers() EmployeeHandlers { return EmployeeHandlers{} }

type employeeData struct {
	UserId     uuid.UUID `json: "user_id"`
	PositionId uuid.UUID `json: "position_id"`
}
