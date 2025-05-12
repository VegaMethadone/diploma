package companylogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/company"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (c CompanyLogic) UpdateCompany(comp *company.Company, companyId, userId uuid.UUID) error {
	// 1. Проверка валидности входных параметров
	if comp == nil {
		logger.NewWarnMessage("Nil company pointer",
			zap.String("operation", "UpdateCompany"),
			zap.Time("time", time.Now()),
		)
		return errors.New("company data cannot be nil")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID",
			zap.String("operation", "UpdateCompany"),
			zap.Time("time", time.Now()),
		)
		return errors.New("company ID cannot be empty")
	}

	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID",
			zap.String("operation", "UpdateCompany"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee ID cannot be empty")
	}

	// 2. Проверка соответствия ID компании
	if comp.ID != companyId {
		logger.NewWarnMessage("Company ID mismatch",
			zap.String("operation", "UpdateCompany"),
			zap.String("requested_company_id", companyId.String()),
			zap.String("provided_company_id", comp.ID.String()),
		)
		return errors.New("requested company ID doesn't match company data")
	}

	// 3. Инициализация подключения к базе данных
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "UpdateCompany"),
			zap.String("company_id", companyId.String()),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 4. Настройка контекста с таймаутом 10 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 5. Начало транзакции (режим read-write для операции обновления)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "UpdateCompany"),
			zap.String("company_id", companyId.String()),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// Откат транзакции при возникновении ошибки
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "UpdateCompany"),
					zap.String("company_id", companyId.String()),
				)
			}
		}
	}()

	ps := postgres.NewPostgresDB()
	// 7. Обновление данных компании
	err = ps.Company.UpdateCompany(ctx, tx, comp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Company not found",
				zap.String("company_id", companyId.String()),
			)
			return fmt.Errorf("company not found: %w", err)
		}

		logger.NewErrMessage("Company update failed",
			zap.Error(err),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("company update failed: %w", err)
	}

	// 8. Фиксация транзакции
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "UpdateCompany"),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	// 9. Логирование успешного обновления
	logger.NewInfoMessage("Company data updated successfully",
		zap.String("company_id", companyId.String()),
		zap.String("updated_by", userId.String()),
		zap.Time("updated_at", time.Now()),
	)

	return nil
}
