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

func NewDepPosition(
	generatedId,
	departmentId uuid.UUID,
	lvl int,
	name string,
) *DepPosition {
	return &DepPosition{
		Id:           generatedId,
		DepartmentId: departmentId,
		Level:        lvl,
		Name:         name,
	}
}
