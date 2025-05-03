package employeelogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/employee"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (e EmployeeLogic) GetEmployee(userId, companyId uuid.UUID) (*employee.Employee, error) {
	// 1. Валидация входных параметров
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user id provided",
			zap.String("operation", "GetEmployee"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("user id cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company id provided",
			zap.String("operation", "GetEmployee"),
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
			zap.String("operation", "GetEmployee"),
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
			zap.String("operation", "GetEmployee"),
			zap.String("user_id", userId.String()),
			zap.String("company_id", companyId.String()),
		)
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}
	defer tx.Rollback()

	// 6. Получение данных сотрудника
	ps := postgres.NewPostgresDB()
	employee, err := ps.Employee.GetEmployeeByUserId(ctx, tx, userId, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Employee not found",
				zap.String("user_id", userId.String()),
				zap.String("company_id", companyId.String()),
			)
			return nil, fmt.Errorf("employee not found")
		}

		logger.NewErrMessage("Failed to fetch employee",
			zap.Error(err),
			zap.String("user_id", userId.String()),
			zap.String("company_id", companyId.String()),
		)
		return nil, fmt.Errorf("failed to fetch employee: %w", err)
	}

	// 7. Проверка принадлежности сотрудника к компании
	if employee.CompanyID != companyId {
		logger.NewWarnMessage("Employee doesn't belong to company",
			zap.String("user_id", userId.String()),
			zap.String("employee_company_id", employee.CompanyID.String()),
			zap.String("requested_company_id", companyId.String()),
		)
		return nil, fmt.Errorf("employee doesn't belong to requested company")
	}

	// 8. Логирование успешного выполнения
	logger.NewInfoMessage("Employee retrieved successfully",
		zap.String("employee_id", employee.ID.String()),
		zap.String("user_id", userId.String()),
		zap.String("company_id", companyId.String()),
		zap.Time("access_time", time.Now()),
	)

	return employee, nil
}
