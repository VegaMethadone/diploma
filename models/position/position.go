package position

import "github.com/google/uuid"

type Position struct {
	Id        uuid.UUID
	CompanyId uuid.UUID
	Lvl       int
	Name      string
}

func NewPosition(newId, companyId uuid.UUID, lvl int, name string) *Position {
	return &Position{
		Id:        newId,
		CompanyId: companyId,
		Lvl:       lvl,
		Name:      name,
	}
}
