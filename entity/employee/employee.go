package employee

import "github.com/google/uuid"

type Employee struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	CompanyId  uuid.UUID
	PositionId uuid.UUID
}

func NewEmployee(userId, companyId, positionId uuid.UUID) *Employee {
	return &Employee{
		Id:         uuid.New(),
		UserId:     userId,
		CompanyId:  companyId,
		PositionId: positionId,
	}
}
