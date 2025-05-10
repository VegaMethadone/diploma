package depemployeeposlogic

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

func (d DepemploeePosLogic) DeleteDepEmployeePos(
	currentlvl int,
	employeeId,
	departmentId,
	positionId uuid.UUID,
) error {
	// 1. Validate input parameters
	if employeeId == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID provided",
			zap.String("operation", "DeleteDepEmployeePos"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee ID cannot be empty")
	}

	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "DeleteDepEmployeePos"),
			zap.Time("time", time.Now()),
		)
		return errors.New("department ID cannot be empty")
	}

	if positionId == uuid.Nil {
		logger.NewWarnMessage("Empty position ID provided",
			zap.String("operation", "DeleteDepEmployeePos"),
			zap.Time("time", time.Now()),
		)
		return errors.New("position ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "DeleteDepEmployeePos"),
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
			zap.String("operation", "DeleteDepEmployeePos"),
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
					zap.String("operation", "DeleteDepEmployeePos"),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "DeleteDepEmployeePos"),
			)
		}
	}()

	ps := postgres.NewPostgresDB()

	// 6. Check if current user has admin rights (level <= 1)
	if currentlvl <= 1 {
		err = ps.DepartmentEmployeePosition.DeleteDepartmentPosition(ctx, tx, positionId)
		if err != nil {
			logger.NewErrMessage("Failed to delete department position",
				zap.Error(err),
				zap.String("operation", "DeleteDepEmployeePos"),
				zap.String("position_id", positionId.String()),
			)
			return fmt.Errorf("failed to delete department position: %w", err)
		}

		logger.NewInfoMessage("Successfully deleted department position (admin override)",
			zap.String("operation", "DeleteDepEmployeePos"),
			zap.String("position_id", positionId.String()),
			zap.String("deleted_by_employee_id", employeeId.String()),
		)
		return nil
	}

	// 7. For non-admin users - verify permissions
	fetchedEmployeeDep, err := ps.DepartmentEmployee.GetEmployeeDepartmentByEmployeeId(ctx, tx, employeeId, departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to verify employee department",
			zap.Error(err),
			zap.String("operation", "DeleteDepEmployeePos"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("failed to verify employee department: %w", err)
	}

	if fetchedEmployeeDep == nil {
		logger.NewWarnMessage("Employee not found in department",
			zap.String("operation", "DeleteDepEmployeePos"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return errors.New("employee not found in specified department")
	}

	// 8. Check if employee has sufficient privileges
	fetchedPos, err := ps.DepartmentEmployeePosition.GetDepartmentPositionById(ctx, tx, fetchedEmployeeDep.PositionID)
	if err != nil {
		logger.NewErrMessage("Failed to verify employee position",
			zap.Error(err),
			zap.String("operation", "DeleteDepEmployeePos"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("failed to verify employee position: %w", err)
	}

	if fetchedPos.Level <= 1 {
		err = ps.DepartmentEmployeePosition.DeleteDepartmentPosition(ctx, tx, positionId)
		if err != nil {
			logger.NewErrMessage("Failed to delete department position",
				zap.Error(err),
				zap.String("operation", "DeleteDepEmployeePos"),
				zap.String("position_id", positionId.String()),
			)
			return fmt.Errorf("failed to delete department position: %w", err)
		}

		logger.NewInfoMessage("Successfully deleted department position",
			zap.String("operation", "DeleteDepEmployeePos"),
			zap.String("position_id", positionId.String()),
			zap.String("deleted_by_employee_id", employeeId.String()),
		)
		return nil
	}

	logger.NewWarnMessage("Access denied for position deletion",
		zap.String("operation", "DeleteDepEmployeePos"),
		zap.String("employee_id", employeeId.String()),
		zap.String("position_id", positionId.String()),
		zap.Int("employee_level", fetchedPos.Level),
	)
	return errors.New("access denied: insufficient privileges")
}
