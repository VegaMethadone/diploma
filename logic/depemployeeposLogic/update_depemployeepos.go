package depemployeeposlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/depposition"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (d DepemploeePosLogic) UpdateDepEmployeePos(
	currentlvl int,
	employeeId,
	departmentId uuid.UUID,
	position *depposition.DepPosition,
) error {
	// 1. Validate input parameters
	if employeeId == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID provided",
			zap.String("operation", "UpdateDepEmployeePos"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee ID cannot be empty")
	}

	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "UpdateDepEmployeePos"),
			zap.Time("time", time.Now()),
		)
		return errors.New("department ID cannot be empty")
	}

	if position == nil {
		logger.NewWarnMessage("Nil position provided",
			zap.String("operation", "UpdateDepEmployeePos"),
		)
		return errors.New("position cannot be nil")
	}

	if position.Id == uuid.Nil {
		logger.NewWarnMessage("Empty position ID provided",
			zap.String("operation", "UpdateDepEmployeePos"),
		)
		return errors.New("position ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "UpdateDepEmployeePos"),
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
			zap.String("operation", "UpdateDepEmployeePos"),
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
					zap.String("operation", "UpdateDepEmployeePos"),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "UpdateDepEmployeePos"),
			)
		}
	}()

	ps := postgres.NewPostgresDB()

	// 6. Check if position belongs to department
	existingPos, err := ps.DepartmentEmployeePosition.GetDepartmentPositionById(ctx, tx, position.Id)
	if err != nil {
		logger.NewErrMessage("Failed to verify position",
			zap.Error(err),
			zap.String("operation", "UpdateDepEmployeePos"),
			zap.String("position_id", position.Id.String()),
		)
		return fmt.Errorf("failed to verify position: %w", err)
	}

	if existingPos.DepartmentId != departmentId {
		logger.NewWarnMessage("Position doesn't belong to department",
			zap.String("operation", "UpdateDepEmployeePos"),
			zap.String("position_department_id", existingPos.DepartmentId.String()),
			zap.String("expected_department_id", departmentId.String()),
		)
		return errors.New("position doesn't belong to specified department")
	}

	// 7. Check if current user has admin rights (level <= 1)
	if currentlvl <= 1 {
		err = ps.DepartmentEmployeePosition.UpdateDepartmentPosition(ctx, tx, position)
		if err != nil {
			logger.NewErrMessage("Failed to update department position",
				zap.Error(err),
				zap.String("operation", "UpdateDepEmployeePos"),
				zap.String("position_id", position.Id.String()),
			)
			return fmt.Errorf("failed to update department position: %w", err)
		}

		logger.NewInfoMessage("Successfully updated department position (admin override)",
			zap.String("operation", "UpdateDepEmployeePos"),
			zap.String("position_id", position.Id.String()),
			zap.String("updated_by_employee_id", employeeId.String()),
		)
		return nil
	}

	// 8. For non-admin users - verify permissions
	fetchedEmployeeDep, err := ps.DepartmentEmployee.GetEmployeeDepartmentByEmployeeId(ctx, tx, employeeId, departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to verify employee department",
			zap.Error(err),
			zap.String("operation", "UpdateDepEmployeePos"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("failed to verify employee department: %w", err)
	}

	if fetchedEmployeeDep == nil {
		logger.NewWarnMessage("Employee not found in department",
			zap.String("operation", "UpdateDepEmployeePos"),
			zap.String("employee_id", employeeId.String()),
			zap.String("department_id", departmentId.String()),
		)
		return errors.New("employee not found in specified department")
	}

	// 9. Check if employee has sufficient privileges
	employeePos, err := ps.DepartmentEmployeePosition.GetDepartmentPositionById(ctx, tx, fetchedEmployeeDep.PositionID)
	if err != nil {
		logger.NewErrMessage("Failed to verify employee position",
			zap.Error(err),
			zap.String("operation", "UpdateDepEmployeePos"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("failed to verify employee position: %w", err)
	}

	if employeePos.Level <= 1 {
		err = ps.DepartmentEmployeePosition.UpdateDepartmentPosition(ctx, tx, position)
		if err != nil {
			logger.NewErrMessage("Failed to update department position",
				zap.Error(err),
				zap.String("operation", "UpdateDepEmployeePos"),
				zap.String("position_id", position.Id.String()),
			)
			return fmt.Errorf("failed to update department position: %w", err)
		}

		logger.NewInfoMessage("Successfully updated department position",
			zap.String("operation", "UpdateDepEmployeePos"),
			zap.String("position_id", position.Id.String()),
			zap.String("updated_by_employee_id", employeeId.String()),
		)
		return nil
	}

	logger.NewWarnMessage("Access denied for position update",
		zap.String("operation", "UpdateDepEmployeePos"),
		zap.String("employee_id", employeeId.String()),
		zap.String("position_id", position.Id.String()),
		zap.Int("employee_level", employeePos.Level),
	)
	return errors.New("access denied: insufficient privileges")
}
