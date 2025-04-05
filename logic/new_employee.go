package logic

import (
	"errors"
	pscompany "labyrinth/database/postgres/pscompany"
	psemployee "labyrinth/database/postgres/psemployee"
	psuser "labyrinth/database/postgres/psuser"

	"github.com/google/uuid"
)

func NewEmployee(userId, companyId, positionId uuid.UUID) error {
	hasUserId, err := psuser.CheckUser(userId)
	if err != nil {
		return err
	}

	hasCompanyId, err := pscompany.CheckCompany(companyId)
	if err != nil {
		return err
	}

	if hasUserId && hasCompanyId {
		return psemployee.NewEmployee(userId, companyId, positionId)
	}

	return errors.New("user or company not exists") // добавить логирование, чтобы было понятно какая компания или юзер не сущесвует
}
