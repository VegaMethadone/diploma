package depemployee

import "github.com/google/uuid"

type DepEmployee struct {
	Id           uuid.UUID `json: "id"`
	EmployeeId   uuid.UUID `json: "employeeId"`
	DepartmentId uuid.UUID `json: "departmentId"`
	AccessId     uuid.UUID `json: "accessId"`
}

func NewDepEmployee(newId, employeeId, departmentId, accessId uuid.UUID) *DepEmployee {
	return &DepEmployee{
		Id:           newId,
		EmployeeId:   employeeId,
		DepartmentId: departmentId,
		AccessId:     accessId,
	}
}
