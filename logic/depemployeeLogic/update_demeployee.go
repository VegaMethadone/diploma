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

func (d DepemployeeLogic) UpdateDepEmployee(
	employeeId,
	departmentId uuid.UUID,
	updatedDepEmployee *depemployee.DepartmentEmployee,
) error {
	// 1. Validate input parameters
	if employeeId == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID provided",
			zap.String("operation", "UpdateDepEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee ID cannot be empty")
	}

	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "UpdateDepEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("department ID cannot be empty")
	}

	if updatedDepEmployee == nil {
		logger.NewWarnMessage("Empty department employee data provided",
			zap.String("operation", "UpdateDepEmployee"),
		)
		return errors.New("department employee data cannot be nil")
	}

	// 2. Validate consistency of IDs
	if updatedDepEmployee.DepartmentID != departmentId {
		logger.NewWarnMessage("Department ID mismatch",
			zap.String("operation", "UpdateDepEmployee"),
			zap.String("expected_department_id", departmentId.String()),
			zap.String("provided_department_id", updatedDepEmployee.DepartmentID.String()),
		)
		return errors.New("department ID in payload doesn't match path parameter")
	}

	if updatedDepEmployee.EmployeeID != employeeId {
		logger.NewWarnMessage("Employee ID mismatch",
			zap.String("operation", "UpdateDepEmployee"),
			zap.String("expected_employee_id", employeeId.String()),
			zap.String("provided_employee_id", updatedDepEmployee.EmployeeID.String()),
		)
		return errors.New("employee ID in payload doesn't match path parameter")
	}

	// 3. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "UpdateDepEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 4. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 5. Begin transaction
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "UpdateDepEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// 6. Ensure proper transaction handling
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "UpdateDepEmployee"),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "UpdateDepEmployee"),
			)
		}
	}()

	// 7. Verify position belongs to department
	ps := postgres.NewPostgresDB()
	fetchedPosition, err := ps.DepartmentEmployeePosition.GetDepartmentPositionById(ctx, tx, updatedDepEmployee.PositionID)
	if err != nil {
		logger.NewErrMessage("Failed to verify position",
			zap.Error(err),
			zap.String("operation", "UpdateDepEmployee"),
			zap.String("position_id", updatedDepEmployee.PositionID.String()),
		)
		return fmt.Errorf("failed to verify position: %w", err)
	}

	if fetchedPosition.DepartmentId != departmentId {
		logger.NewWarnMessage("Position doesn't belong to department",
			zap.String("operation", "UpdateDepEmployee"),
			zap.String("position_department_id", fetchedPosition.DepartmentId.String()),
			zap.String("expected_department_id", departmentId.String()),
		)
		return errors.New("position doesn't belong to specified department")
	}

	// 8. Update department employee
	err = ps.DepartmentEmployee.UpdateEmployeeDepartment(ctx, tx, updatedDepEmployee)
	if err != nil {
		logger.NewErrMessage("Failed to update department employee",
			zap.Error(err),
			zap.String("operation", "UpdateDepEmployee"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("failed to update department employee: %w", err)
	}

	logger.NewInfoMessage("Successfully updated department employee",
		zap.String("operation", "UpdateDepEmployee"),
		zap.String("employee_id", employeeId.String()),
		zap.String("department_id", departmentId.String()),
		zap.String("position_id", updatedDepEmployee.PositionID.String()),
	)

	return nil
}
