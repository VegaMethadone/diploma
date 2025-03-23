package logic

import (
	"labyrinth/entity/company"

	"github.com/google/uuid"
)

func NewCompany(name, description string, owner uuid.UUID) error {
	company_ := company.NewCompany(owner, name, description)

	return ps.RegisterCompany(company_)
}
