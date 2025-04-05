package department

import "github.com/google/uuid"

type Department struct {
	Id          uuid.UUID `json: "id"`
	CompanyId   uuid.UUID `json: "companyId"`
	OwnerId     uuid.UUID `json: "ownerId"`
	Name        string    `json: "name"`
	Description string    `json: "description"`
	Parent      uuid.UUID `json: "parent"`
}

func NewDepartment(companyId, employeeId, parent uuid.UUID, name, descritpion string) *Department {
	return &Department{
		Id:          uuid.New(),
		CompanyId:   companyId,
		OwnerId:     employeeId,
		Name:        name,
		Description: descritpion,
		Parent:      parent,
	}
}
