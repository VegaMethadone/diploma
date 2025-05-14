package depposition

import "labyrinth/logic"

var bl *logic.BusinessLogic = logic.NewBusinessLogic()

type DepPositionHandlers struct{}

func NewDepPositionHandlers() DepPositionHandlers { return DepPositionHandlers{} }

type deppositionData struct {
	Lvl  int    `json: "lvl"`
	Name string `json: "name"`
}
