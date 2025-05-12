package companylogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/company"
	"labyrinth/models/employee"
	"labyrinth/models/position"
	notebookLogic "labyrinth/notebook/logic"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (c CompanyLogic) NewCompany(userId uuid.UUID, name, description string) (uuid.UUID, error) {
	// 1. Валидация входных данных
	if userId == uuid.Nil {
		return uuid.Nil, errors.New("user id cannot be empty")
	}
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" || description == "" {
		return uuid.Nil, errors.New("name and description cannot be empty")
	}

	// 2. Инициализация подключения к БД
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "NewCompany"),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, fmt.Errorf("database connection failed: %w", err)
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
		return uuid.Nil, fmt.Errorf("transaction begin failed: %w", err)
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
			return uuid.Nil, fmt.Errorf("user not found: %w", err)
		}

		logger.NewErrMessage("Failed to get user",
			zap.Error(err),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 6. Генерация нового UUID для компании
	newCompanyUUID, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate company UUID",
			zap.Error(err),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to generate company UUID: %w", err)
	}

	// 7. Создание новой компании
	newCompany := company.NewCompany(
		userId,
		name,
		description,
		"no address",
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
		return uuid.Nil, fmt.Errorf("failed to create company: %w", err)
	}

	// 8. Добавление пользователя в компанию
	err = ps.Company.AddUserToCompany(ctx, tx, userId, newCompanyUUID)
	if err != nil {
		logger.NewErrMessage("Failed to add user to company",
			zap.Error(err),
			zap.String("user_id", userId.String()),
			zap.String("company_id", newCompanyUUID.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to add user to company: %w", err)
	}

	employeeUUID, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate employee UUID",
			zap.Error(err),
			zap.String("user_id", userId.String()),
			zap.String("company_id", newCompanyUUID.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to generate employee UUID: %w", err)
	}

	positionUUID, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate position UUID",
			zap.Error(err),
			zap.String("user_id", userId.String()),
			zap.String("company_id", newCompanyUUID.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to generate position UUID: %w", err)
	}

	// Создание позиции владельца
	newPosition := position.NewPosition(positionUUID, newCompanyUUID, 0, "owner")
	err = ps.Position.CreatePosition(ctx, tx, &newPosition)
	if err != nil {
		logger.NewErrMessage("Failed to create owner position",
			zap.Error(err),
			zap.String("user_id", userId.String()),
			zap.String("company_id", newCompanyUUID.String()),
			zap.String("position_id", positionUUID.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to create owner position: %w", err)
	}

	// Создание сотрудника (владельца)
	newEmployee := employee.NewEmployee(employeeUUID, userId, newCompanyUUID, positionUUID)
	err = ps.Employee.CreateEmployee(ctx, tx, newEmployee)
	if err != nil {
		logger.NewErrMessage("Failed to create employee record",
			zap.Error(err),
			zap.String("user_id", userId.String()),
			zap.String("company_id", newCompanyUUID.String()),
			zap.String("employee_id", employeeUUID.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to create employee record: %w", err)
	}

	// Создание корневой папки компании в файловой системе
	fileSystem := notebookLogic.NewFileSystem()
	err = fileSystem.Folder.CreateFolder(
		employeeUUID,
		newCompanyUUID,
		newCompanyUUID,
		newCompanyUUID,
		true,
		name,
		description,
	)
	if err != nil {
		logger.NewErrMessage("Failed to create company root folder",
			zap.Error(err),
			zap.String("user_id", userId.String()),
			zap.String("company_id", newCompanyUUID.String()),
			zap.String("employee_id", employeeUUID.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to create company root folder: %w", err)
	}

	// 10. Коммит транзакции
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "NewCompany"),
			zap.String("company_id", newCompanyUUID.String()),
		)
		return uuid.Nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	// 11. Логирование успешного создания
	logger.NewInfoMessage("Company created successfully",
		zap.String("company_name", name),
		zap.String("user_id", userId.String()),
		zap.String("company_id", newCompanyUUID.String()),
		zap.Time("created_at", time.Now()),
	)

	return newCompanyUUID, nil
}
