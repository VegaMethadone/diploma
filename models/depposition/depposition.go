package depposition

import (
	"github.com/google/uuid"
)

type DepPosition struct {
	Id           uuid.UUID `json: "id"`
	DepartmentId uuid.UUID `json: "departmentId"`
	Level        int       `json: "lvl"`
	Name         string    `json: "name"`
}
