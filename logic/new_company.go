package logic

import (
	"fmt"
	"labyrinth/entity/company"
	"labyrinth/entity/position"

	"github.com/google/uuid"
)

func NewCompany(name, description, positionName string, owner uuid.UUID) error {
	// Создаем новую компанию
	company_ := company.NewCompany(owner, name, description)

	// Регистрируем компанию в хранилище
	companyId, err := ps.RegisterCompany(company_)
	if err != nil {
		return fmt.Errorf("failed to add new company: %w", err)
	}

	// Создаем позицию по умолчанию
	pos := position.NewPosition(*companyId, 0, positionName)
	err = ps.NewPosition(pos)
	if err != nil {
		return fmt.Errorf("failed to add new position: %w", err)
	}

	// Добавляем владельца как сотрудника
	err = ps.NewEmployee(owner, *companyId, pos.Id)
	if err != nil {
		return fmt.Errorf("failed to add new employee: %w", err)
	}

	return nil
}
