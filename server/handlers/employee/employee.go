package employee

import (
	"labyrinth/logic"

	"github.com/google/uuid"
)

var bl *logic.BusinessLogic = logic.NewBusinessLogic()

type employeeData struct {
	UserId     uuid.UUID `json: "user_id"`
	PositionId uuid.UUID `json: "position_id"`
}
