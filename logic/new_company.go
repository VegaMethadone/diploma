package logic

import (
	"fmt"
	pscompany "labyrinth/database/postgres/pscompany"
	psemployee "labyrinth/database/postgres/psemployee"
	psposition "labyrinth/database/postgres/psposition"
	"labyrinth/entity/company"
	"labyrinth/entity/position"

	"github.com/google/uuid"
)

func NewCompany(name, description, positionName string, owner uuid.UUID) error {
	// Создаем новую компанию
	company_ := company.NewCompany(owner, name, description)

	// Регистрируем компанию в хранилище
	companyId, err := pscompany.RegisterCompany(company_)
	if err != nil {
		return fmt.Errorf("failed to add new company: %w", err)
	}

	// Создаем позицию по умолчанию
	pos := position.NewPosition(*companyId, 0, positionName)
	err = psposition.NewPosition(pos)
	if err != nil {
		return fmt.Errorf("failed to add new position: %w", err)
	}

	// Добавляем владельца как сотрудника
	err = psemployee.NewEmployee(owner, *companyId, pos.Id)
	if err != nil {
		return fmt.Errorf("failed to add new employee: %w", err)
	}

	return nil
}
