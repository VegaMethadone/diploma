package depemployeelogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/depemployee"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (d DepemployeeLogic) NewDepemployee(
	employeeId,
	departmentId,
	positionId uuid.UUID,
) error {
	// 1. Validate all input parameters
	if employeeId == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID provided",
			zap.String("operation", "NewDepemployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee ID cannot be empty")
	}

	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "NewDepemployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("department ID cannot be empty")
	}

	if positionId == uuid.Nil {
		logger.NewWarnMessage("Empty position ID provided",
			zap.String("operation", "NewDepemployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("position ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "NewDepemployee"),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin transaction
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "NewDepemployee"),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// Ensure transaction is rolled back on error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "NewDepemployee"),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "NewDepemployee"),
			)
		}
	}()

	// 5. Generate and validate new UUID
	ps := postgres.NewPostgresDB()
	generatedId, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate UUID",
			zap.Error(err),
			zap.String("operation", "NewDepemployee"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	// 6. Create new department employee
	newDepEmployee := depemployee.NewDepartmentEmployee(
		generatedId,
		employeeId,
		departmentId,
		positionId,
	)

	// 7. Save to database
	err = ps.DepartmentEmployee.CreateEmployeeDepartment(ctx, tx, newDepEmployee)
	if err != nil {
		logger.NewErrMessage("Failed to create department employee",
			zap.Error(err),
			zap.String("operation", "NewDepemployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
			zap.String("position_id", positionId.String()),
		)
		return fmt.Errorf("failed to create department employee: %w", err)
	}

	logger.NewInfoMessage("Successfully created department employee",
		zap.String("operation", "NewDepemployee"),
		zap.String("employee_id", employeeId.String()),
		zap.String("department_id", departmentId.String()),
		zap.String("position_id", positionId.String()),
		zap.String("new_depemployee_id", generatedId.String()),
	)

	return nil
}
