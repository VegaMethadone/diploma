package companylogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/company"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// func (c CompanyLogic) NewCompany(userId uuid.UUID, name, description string) error {
// 	// 1. Валидация входных данных
// 	if userId == uuid.Nil {
// 		return errors.New("user id cannot be empty")
// 	}
// 	if strings.TrimSpace(name) == "" || strings.TrimSpace(description) == "" {
// 		return errors.New("name or description cannot be empty")
// 	}

// 	// 2. Инициализация подключения к БД
// 	db, err := sql.Open("postgres", postgres.GetConnection())
// 	if err != nil {
// 		logger.NewErrMessage("Database connection failed",
// 			zap.Error(err),
// 			zap.String("operation", "NewCompany"),
// 			zap.String("user_id", userId.String()),
// 		)
// 		return fmt.Errorf("database connection failed: %w", err)
// 	}
// 	defer db.Close()

// 	// 3. Настройка контекста с таймаутом
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// 4. Начало транзакции
// 	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
// 	if err != nil {
// 		logger.NewErrMessage("Transaction begin failed",
// 			zap.Error(err),
// 			zap.String("operation", "NewCompany"),
// 			zap.String("user_id", userId.String()),
// 		)
// 		return fmt.Errorf("transaction begin failed: %w", err)
// 	}

// 	// Обеспечиваем откат транзакции в случае ошибки
// 	defer func() {
// 		if err != nil {
// 			if rbErr := tx.Rollback(); rbErr != nil {
// 				logger.NewErrMessage("Transaction rollback failed",
// 					zap.Error(rbErr),
// 					zap.String("operation", "NewCompany"),
// 				)
// 			}
// 		}
// 	}()

// 	// 5. Получение данных пользователя
// 	ps := postgres.NewPostgresDB()
// 	foundUser, err := ps.User.GetUserByID(ctx, tx, userId)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			logger.NewWarnMessage("User not found",
// 				zap.String("user_id", userId.String()),
// 			)
// 			return fmt.Errorf("user not found: %w", err)
// 		}

// 		logger.NewErrMessage("Failed to get user",
// 			zap.Error(err),
// 			zap.String("user_id", userId.String()),
// 		)
// 		return fmt.Errorf("failed to get user: %w", err)
// 	}

// 	newCompanyUUID, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
// 	if err != nil {

// 	}
// 	// 6. Создание новой компании
// 	newCompany := company.NewCompany(
// 		userId,
// 		name,
// 		description,
// 		"no address", // временное значение, можно вынести в конфиг
// 		foundUser.Phone,
// 		foundUser.Email,
// 	)
// 	newCompany.ID = newCompanyUUID

// 	err = ps.Company.CreateCompany(ctx, tx, newCompany)
// 	if err != nil {
// 		logger.NewErrMessage("Failed to create company",
// 			zap.Error(err),
// 			zap.String("company_name", name),
// 			zap.String("user_id", userId.String()),
// 		)
// 		return fmt.Errorf("failed to create company: %w", err)
// 	}

// 	err = ps.Company.AddUserToCompany(ctx, tx, userId, newCompanyUUID)
// 	if err != nil {

// 	}

// 	// добавить minio

// 	// 7. Коммит транзакции
// 	if err = tx.Commit(); err != nil {
// 		logger.NewErrMessage("Transaction commit failed",
// 			zap.Error(err),
// 			zap.String("operation", "NewCompany"),
// 		)
// 		return fmt.Errorf("transaction commit failed: %w", err)
// 	}

// 	// 8. Логирование успешного создания
// 	logger.NewInfoMessage("Company created successfully",
// 		zap.String("company_name", name),
// 		zap.String("user_id", userId.String()),
// 		zap.Time("created_at", time.Now()),
// 	)

// 	return nil
// }

func (c CompanyLogic) NewCompany(userId uuid.UUID, name, description string) error {
	// 1. Валидация входных данных
	if userId == uuid.Nil {
		return errors.New("user id cannot be empty")
	}
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" || description == "" {
		return errors.New("name and description cannot be empty")
	}

	// 2. Инициализация подключения к БД
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "NewCompany"),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Настройка контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Начало транзакции
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "NewCompany"),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// Обеспечиваем откат транзакции в случае ошибки
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "NewCompany"),
					zap.String("user_id", userId.String()),
				)
			}
		}
	}()

	// 5. Получение данных пользователя
	ps := postgres.NewPostgresDB()
	foundUser, err := ps.User.GetUserByID(ctx, tx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("User not found",
				zap.String("user_id", userId.String()),
			)
			return fmt.Errorf("user not found: %w", err)
		}

		logger.NewErrMessage("Failed to get user",
			zap.Error(err),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 6. Генерация нового UUID для компании
	newCompanyUUID, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate company UUID",
			zap.Error(err),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("failed to generate company UUID: %w", err)
	}

	// 7. Создание новой компании
	newCompany := company.NewCompany(
		userId,
		name,
		description,
		"no address", // временное значение, можно вынести в конфиг
		foundUser.Phone,
		foundUser.Email,
	)
	newCompany.ID = newCompanyUUID

	err = ps.Company.CreateCompany(ctx, tx, newCompany)
	if err != nil {
		logger.NewErrMessage("Failed to create company",
			zap.Error(err),
			zap.String("company_name", name),
			zap.String("user_id", userId.String()),
			zap.String("company_id", newCompanyUUID.String()),
		)
		return fmt.Errorf("failed to create company: %w", err)
	}

	// 8. Добавление пользователя в компанию
	err = ps.Company.AddUserToCompany(ctx, tx, userId, newCompanyUUID)
	if err != nil {
		logger.NewErrMessage("Failed to add user to company",
			zap.Error(err),
			zap.String("user_id", userId.String()),
			zap.String("company_id", newCompanyUUID.String()),
		)
		return fmt.Errorf("failed to add user to company: %w", err)
	}

	// 9. Сохранение данных в MinIO (если требуется)
	// Пример:
	/*
	   if err := c.saveCompanyLogoToMinio(ctx, newCompanyUUID); err != nil {
	       logger.NewWarnMessage("Failed to save company logo to MinIO",
	           zap.Error(err),
	           zap.String("company_id", newCompanyUUID.String()),
	       )
	       // Можно продолжить, если это не критическая ошибка
	   }
	*/

	// 10. Коммит транзакции
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "NewCompany"),
			zap.String("company_id", newCompanyUUID.String()),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	// 11. Логирование успешного создания
	logger.NewInfoMessage("Company created successfully",
		zap.String("company_name", name),
		zap.String("user_id", userId.String()),
		zap.String("company_id", newCompanyUUID.String()),
		zap.Time("created_at", time.Now()),
	)

	return nil
}
