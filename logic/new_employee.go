package logic

import (
	"errors"

	"github.com/google/uuid"
)

func NewEmployee(userId, companyId uuid.UUID) error {
	hasUserId, err := ps.CheckUser(userId)
	if err != nil {
		return err
	}

	hasCompanyId, err := ps.CheckCompany(companyId)
	if err != nil {
		return err
	}

	if hasUserId && hasCompanyId {
		return ps.NewEmployee(userId, companyId)
	}

	return errors.New("user or company not exists") // добавить логирование, чтобы было понятно какая компания или юзер не сущесвует
}
