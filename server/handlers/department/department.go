package department

import (
	"labyrinth/logic"
	notebookLogic "labyrinth/notebook/logic"

	"github.com/google/uuid"
)

type DepartmentHandlers struct{}

func NewDepartmentHandlers() DepartmentHandlers { return DepartmentHandlers{} }

var bl *logic.BusinessLogic = logic.NewBusinessLogic()
var fsl *notebookLogic.FileSystem = notebookLogic.NewFileSystem()

type departmentData struct {
	ParentId    uuid.UUID `json: "id"`
	Name        string    `json: "name"`
	Description string    `json: "description"`
}
