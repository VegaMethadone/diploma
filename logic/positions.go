package logic

import (
	"encoding/json"
	"fmt"
	psposition "labyrinth/database/postgres/psposition"
	"labyrinth/entity/employee"
	"labyrinth/entity/position"

	"github.com/google/uuid"
)

func AddPosition(whoPositionId uuid.UUID, what *position.Position) error {
	// 1. Проверяем существование позиции, от имени которой добавляем
	currentPos, err := psposition.CheckPosition(whoPositionId)
	if err != nil {
		return fmt.Errorf("failed to verify current position: %w", err)
	}
	if currentPos == nil {
		return fmt.Errorf("current position not found")
	}

	// 2. Проверяем права доступа (иерархия уровней)
	err = verifyLvlAccess(currentPos, what)
	if err != nil {
		return err
	}

	// 3. Создаем новую позицию
	if err := psposition.NewPosition(what); err != nil {
		return fmt.Errorf("failed to create new position: %w", err)
	}

	return nil
}

func DeletePosition(whoPositionId uuid.UUID, what *position.Position) error {
	// 1. Проверяем существование позиции, от имени которой добавляем
	currentPos, err := psposition.CheckPosition(whoPositionId)
	if err != nil {
		return fmt.Errorf("failed to verify current position: %w", err)
	}
	if currentPos == nil {
		return fmt.Errorf("current position not found")
	}

	// 2. Проверяем права доступа (иерархия уровней)
	err = verifyLvlAccess(currentPos, what)
	if err != nil {
		return err
	}

	err = psposition.DeletePosition(what.Id)
	if err != nil {
		return err
	}

	return nil
}

func RenamePosition(whoPositionId uuid.UUID, what *position.Position) error {
	// 1. Проверяем существование позиции, от имени которой добавляем
	currentPos, err := psposition.CheckPosition(whoPositionId)
	if err != nil {
		return fmt.Errorf("failed to verify current position: %w", err)
	}
	if currentPos == nil {
		return fmt.Errorf("current position not found")
	}

	whatPositionCnage, err := psposition.CheckPosition(what.Id)
	if err != nil {
		return fmt.Errorf("failed to verify position that has to be changed: %w", err)
	}
	if currentPos == nil {
		return fmt.Errorf("current position not found")
	}

	// 2. Проверяем права доступа (иерархия уровней)
	err = verifyLvlAccess(currentPos, whatPositionCnage)
	if err != nil {
		return err
	}

	err = psposition.ChangePosition(what)
	if err != nil {
		return err
	}

	return nil
}

func GetPositions(who *employee.Employee) (string, error) {
	positions, err := psposition.GetPositions(who.CompanyId)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(positions)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil

}

func verifyLvlAccess(currentPos, what *position.Position) error {
	if currentPos.Lvl < what.Lvl {
		return fmt.Errorf(
			"permission denied: '%s' (level %d) cannot create '%s' (level %d)",
			currentPos.Name,
			currentPos.Lvl,
			what.Name,
			what.Lvl,
		)
	}
	return nil
}
