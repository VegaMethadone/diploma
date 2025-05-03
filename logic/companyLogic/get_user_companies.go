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

func (c CompanyLogic) GetUserCompanies(userId uuid.UUID) (*[]company.Company, error) {
	// 1. Валидация входных данных
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user id provided",
			zap.String("operation", "GetUserCompanies"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("user id cannot be empty")
	}

	// 2. Инициализация подключения к БД
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "GetUserCompanies"),
			zap.String("user_id", userId.String()),
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
			zap.String("operation", "GetUserCompanies"),
			zap.String("user_id", userId.String()),
		)
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}
	defer tx.Rollback() // Безопасный откат для read-only транзакции

	// 5. Получение компаний пользователя
	ps := postgres.NewPostgresDB()
	companies, err := ps.Company.GetCompaniesByUser(ctx, tx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewInfoMessage("No companies found for user",
				zap.String("user_id", userId.String()),
				zap.Time("time", time.Now()),
			)
			return &[]company.Company{}, nil
		}

		logger.NewErrMessage("Failed to fetch user companies",
			zap.Error(err),
			zap.String("user_id", userId.String()),
			zap.Time("time", time.Now()),
		)
		return nil, fmt.Errorf("failed to fetch user companies: %w", err)
	}

	// 6. вытащить данные из minio если они есть

	// 7. Логирование успешного выполнения
	logger.NewInfoMessage("User companies retrieved successfully",
		zap.String("user_id", userId.String()),
		zap.Time("time", time.Now()),
	)

	return companies, nil
}
