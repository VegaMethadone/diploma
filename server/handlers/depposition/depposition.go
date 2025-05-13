package depposition

import "labyrinth/logic"

var bl *logic.BusinessLogic = logic.NewBusinessLogic()

type deppositionData struct {
	Lvl  int    `json: "lvl"`
	Name string `json: "name"`
}
