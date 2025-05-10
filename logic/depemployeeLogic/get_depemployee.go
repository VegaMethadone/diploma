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

func (d DepemployeeLogic) GetDepartmentEmployee(
	employeeId,
	departmentId uuid.UUID,
) (*depemployee.DepartmentEmployee, error) {
	// 1. Validate input parameters
	if employeeId == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID provided",
			zap.String("operation", "GetDepartmentEmployee"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("employee ID cannot be empty")
	}

	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "GetDepartmentEmployee"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("department ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "GetDepartmentEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return nil, fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin read-only transaction
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "GetDepartmentEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}

	// 5. Ensure proper transaction handling
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "GetDepartmentEmployee"),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "GetDepartmentEmployee"),
			)
		}
	}()

	// 6. Fetch department employee
	ps := postgres.NewPostgresDB()
	fetchedEmployee, err := ps.DepartmentEmployee.GetEmployeeDepartmentByEmployeeId(ctx, tx, employeeId, departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to get department employee",
			zap.Error(err),
			zap.String("operation", "GetDepartmentEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return nil, fmt.Errorf("failed to get department employee: %w", err)
	}

	// 7. Check if employee exists
	if fetchedEmployee == nil {
		logger.NewInfoMessage("Department employee not found",
			zap.String("operation", "GetDepartmentEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return nil, nil
	}

	logger.NewInfoMessage("Successfully retrieved department employee",
		zap.String("operation", "GetDepartmentEmployee"),
		zap.String("employee_id", employeeId.String()),
		zap.String("department_id", departmentId.String()),
	)

	return fetchedEmployee, nil
}
