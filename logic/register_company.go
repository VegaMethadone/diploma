package logic

import (
	"labyrinth/entity/company"

	"github.com/google/uuid"
)

func NewCompany(name, description string, owner uuid.UUID) error {
	company_ := company.NewCompany(owner, name, description)

	companyId, err := ps.RegisterCompany(company_)
	if err != nil {
		return err
	}

	return ps.NewEmployee(owner, *companyId)
}
