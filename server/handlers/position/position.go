package position

import "labyrinth/logic"

const (
	userIDKey string = "id"
)

var bl *logic.BusinessLogic = logic.NewBusinessLogic()

type PositionHandlers struct{}

func NewPositionHandlers() PositionHandlers { return PositionHandlers{} }

type positionData struct {
	Lvl  int    `json: "lvl"`
	Name string `json: "name"`
}
