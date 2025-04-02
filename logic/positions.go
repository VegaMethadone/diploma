package logic

import (
	"fmt"
	"labyrinth/entity/position"

	"github.com/google/uuid"
)

func AddPosition(whoPositionId uuid.UUID, what *position.Position) error {
	// 1. Проверяем существование позиции, от имени которой добавляем
	currentPos, err := ps.CheckPosition(whoPositionId)
	if err != nil {
		return fmt.Errorf("failed to verify current position: %w", err)
	}
	if currentPos == nil {
		return fmt.Errorf("current position not found")
	}

	// 2. Проверяем права доступа (иерархия уровней)
	if currentPos.Lvl < what.Lvl {
		return fmt.Errorf(
			"permission denied: '%s' (level %d) cannot create '%s' (level %d)",
			currentPos.Name,
			currentPos.Lvl,
			what.Name,
			what.Lvl,
		)
	}

	// 3. Создаем новую позицию
	if err := ps.NewPosition(what); err != nil {
		return fmt.Errorf("failed to create new position: %w", err)
	}

	return nil
}
func DeletePosition() {}
func RenamePosition() {}
