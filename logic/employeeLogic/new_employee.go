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

func (e EmployeeLogic) NewEmployee(employeeId, userId, companyId, positionId uuid.UUID) error {
	// 1. Input validation
	if employeeId == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID provided",
			zap.String("operation", "NewEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee ID cannot be empty")
	}

	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "NewEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("user ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "NewEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("company ID cannot be empty")
	}

	if positionId == uuid.Nil {
		logger.NewWarnMessage("Empty position ID provided",
			zap.String("operation", "NewEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("position ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "NewEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin transaction (should be read-write for create operation)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "NewEmployee"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// Ensure transaction is rolled back on error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "NewEmployee"),
				)
			}
		}
	}()

	ps := postgres.NewPostgresDB()

	// 5. Verify position exists in company
	fetchedPosition, err := ps.Position.GetPositionById(ctx, tx, positionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Position not found",
				zap.String("position_id", positionId.String()),
				zap.String("company_id", companyId.String()),
			)
			return fmt.Errorf("position not found: %w", err)
		}

		logger.NewErrMessage("Failed to fetch position",
			zap.Error(err),
			zap.String("position_id", positionId.String()),
		)
		return fmt.Errorf("failed to fetch position: %w", err)
	}

	if fetchedPosition.CompanyID != companyId {
		logger.NewWarnMessage("Position doesn't belong to company",
			zap.String("position_id", positionId.String()),
			zap.String("position_company_id", fetchedPosition.CompanyID.String()),
			zap.String("requested_company_id", companyId.String()),
		)
		return errors.New("position doesn't belong to specified company")
	}

	// // 6. Check if employee already exists
	// exists, err := ps.Employee.ExistsEmployee(ctx, tx, employeeId)
	// if err != nil {
	// 	logger.NewErrMessage("Failed to check employee existence",
	// 		zap.Error(err),
	// 		zap.String("employee_id", employeeId.String()),
	// 		zap.String("company_id", companyId.String()),
	// 	)
	// 	return fmt.Errorf("failed to check employee existence: %w", err)
	// }

	// if exists {
	// 	logger.NewWarnMessage("Employee already exists",
	// 		zap.String("employee_id", employeeId.String()),
	// 		zap.String("company_id", companyId.String()),
	// 	)
	// 	return errors.New("employee already exists in this company")
	// }

	// 7. Generate new UUID for employee
	generatedId, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate employee UUID",
			zap.Error(err),
			zap.String("operation", "NewEmployee"),
		)
		return fmt.Errorf("failed to generate employee UUID: %w", err)
	}

	// 8. Create new employee
	newEmployee := employee.NewEmployee(generatedId, userId, companyId, positionId)
	err = ps.Employee.CreateEmployee(ctx, tx, newEmployee)
	if err != nil {
		logger.NewErrMessage("Failed to create employee",
			zap.Error(err),
			zap.String("employee_id", generatedId.String()),
			zap.String("user_id", userId.String()),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("failed to create employee: %w", err)
	}

	// 9. Commit transaction
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "NewEmployee"),
			zap.String("employee_id", generatedId.String()),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	// 10. Log successful creation
	logger.NewInfoMessage("Employee created successfully",
		zap.String("employee_id", generatedId.String()),
		zap.String("user_id", userId.String()),
		zap.String("company_id", companyId.String()),
		zap.String("position_id", positionId.String()),
		zap.Time("created_at", time.Now()),
	)

	return nil
}
