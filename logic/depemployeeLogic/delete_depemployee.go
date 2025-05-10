package depemployeelogic

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

func (d DepemployeeLogic) DeleteDepartmentEmployee(
	employeeId,
	departmentId,
	depemployeeId uuid.UUID,
) error {
	// 1. Validate input parameters
	if employeeId == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID provided",
			zap.String("operation", "DeleteDepartmentEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee ID cannot be empty")
	}

	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "DeleteDepartmentEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("department ID cannot be empty")
	}

	if depemployeeId == uuid.Nil {
		logger.NewWarnMessage("Empty department employee ID provided",
			zap.String("operation", "DeleteDepartmentEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("department employee ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "DeleteDepartmentEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
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
			zap.String("operation", "DeleteDepartmentEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// 5. Ensure proper transaction handling
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "DeleteDepartmentEmployee"),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "DeleteDepartmentEmployee"),
			)
		}
	}()

	// 6. Verify department employee exists and belongs to specified employee/department
	ps := postgres.NewPostgresDB()
	fetchedDepEmployee, err := ps.DepartmentEmployee.GetEmployeeDepartmentByEmployeeId(ctx, tx, employeeId, departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to verify department employee",
			zap.Error(err),
			zap.String("operation", "DeleteDepartmentEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("failed to verify department employee: %w", err)
	}

	if fetchedDepEmployee == nil {
		logger.NewWarnMessage("Department employee not found",
			zap.String("operation", "DeleteDepartmentEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return errors.New("department employee not found")
	}

	if fetchedDepEmployee.ID != depemployeeId {
		logger.NewWarnMessage("Department employee ID mismatch",
			zap.String("operation", "DeleteDepartmentEmployee"),
			zap.String("expected_depemployee_id", depemployeeId.String()),
			zap.String("found_depemployee_id", fetchedDepEmployee.ID.String()),
		)
		return errors.New("department employee ID doesn't match")
	}

	// 7. Delete department employee
	err = ps.DepartmentEmployee.DeleteEmployeeDepartment(ctx, tx, depemployeeId)
	if err != nil {
		logger.NewErrMessage("Failed to delete department employee",
			zap.Error(err),
			zap.String("operation", "DeleteDepartmentEmployee"),
			zap.String("depemployee_id", depemployeeId.String()),
		)
		return fmt.Errorf("failed to delete department employee: %w", err)
	}

	logger.NewInfoMessage("Successfully deleted department employee",
		zap.String("operation", "DeleteDepartmentEmployee"),
		zap.String("employee_id", employeeId.String()),
		zap.String("department_id", departmentId.String()),
		zap.String("depemployee_id", depemployeeId.String()),
	)

	return nil
}
