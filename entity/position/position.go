package position

import "github.com/google/uuid"

type Position struct {
	Id        uuid.UUID
	CompanyId uuid.UUID
	Lvl       int
	Name      string
}

func NewPosition(companyId uuid.UUID, lvl int, name string) *Position {
	return &Position{
		Id:        uuid.New(),
		CompanyId: companyId,
		Lvl:       lvl,
		Name:      name,
	}
}
