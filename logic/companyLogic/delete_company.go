package companylogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (c CompanyLogic) DeleteCompany(companyId uuid.UUID) error {
	// 1. Проверка валидности входных данных
	if companyId == uuid.Nil {
		logger.NewWarnMessage("Пустой ID компании",
			zap.String("operation", "DeleteCompany"),
			zap.Time("time", time.Now()),
		)
		return errors.New("ID компании не может быть пустым")
	}

	// 2. Подключение к базе данных
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Ошибка подключения к БД",
			zap.Error(err),
			zap.String("operation", "DeleteCompany"),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	defer db.Close() // Гарантированное закрытие соединения

	// 3. Создание контекста с таймаутом (10 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Отмена контекста при выходе

	// 4. Начало транзакции (транзакция должна быть read-write для операций удаления)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Ошибка начала транзакции",
			zap.Error(err),
			zap.String("operation", "DeleteCompany"),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	// Откат транзакции в случае ошибки
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Ошибка отката транзакции",
					zap.Error(rbErr),
					zap.String("operation", "DeleteCompany"),
					zap.String("company_id", companyId.String()),
				)
			}
		}
	}()

	ps := postgres.NewPostgresDB()

	// 5. Удаление компании
	err = ps.Company.DeleteCompany(ctx, tx, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Компания не найдена",
				zap.String("company_id", companyId.String()),
			)
			return fmt.Errorf("компания не найдена: %w", err)
		}

		logger.NewErrMessage("Ошибка удаления компании",
			zap.Error(err),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("ошибка удаления компании: %w", err)
	}

	// 6. Деактивация пользователей компании
	err = ps.Company.DeactivateCompanyUsers(ctx, tx, companyId)
	if err != nil {
		logger.NewErrMessage("Ошибка деактивации пользователей компании",
			zap.Error(err),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("ошибка деактивации пользователей компании: %w", err)
	}

	// 7. Фиксация транзакции
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Ошибка подтверждения транзакции",
			zap.Error(err),
			zap.String("operation", "DeleteCompany"),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	// 8. Логирование успешного удаления
	logger.NewInfoMessage("Компания успешно удалена",
		zap.String("company_id", companyId.String()),
		zap.Time("deleted_at", time.Now()),
	)

	return nil
}
