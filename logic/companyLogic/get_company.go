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

func (c CompanyLogic) GetCompany(userId, companyId uuid.UUID) (*company.Company, error) {
	// 1. Валидация входных параметров
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user id provided",
			zap.String("operation", "GetCompany"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("user id cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company id provided",
			zap.String("operation", "GetCompany"),
			zap.String("user_id", userId.String()),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("company id cannot be empty")
	}

	// 2. Инициализация подключения к БД
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "GetCompany"),
			zap.String("user_id", userId.String()),
			zap.String("company_id", companyId.String()),
		)
		return nil, fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Настройка контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Начало read-only транзакции
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "GetCompany"),
			zap.String("user_id", userId.String()),
			zap.String("company_id", companyId.String()),
		)
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}
	defer tx.Rollback()

	// 5. Получение данных компании
	ps := postgres.NewPostgresDB()
	foundCompany, err := ps.Company.GetCompanyByID(ctx, tx, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Company not found",
				zap.String("company_id", companyId.String()),
				zap.String("user_id", userId.String()),
			)
			return nil, fmt.Errorf("company not found")
		}

		logger.NewErrMessage("Failed to fetch company",
			zap.Error(err),
			zap.String("company_id", companyId.String()),
			zap.String("user_id", userId.String()),
		)
		return nil, fmt.Errorf("failed to fetch company: %w", err)
	}

	// 6. Логирование успешного выполнения
	logger.NewInfoMessage("Company retrieved successfully",
		zap.String("company_id", companyId.String()),
		zap.String("user_id", userId.String()),
		zap.Time("access_time", time.Now()),
	)

	return foundCompany, nil
}
