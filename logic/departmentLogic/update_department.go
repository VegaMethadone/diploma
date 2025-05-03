package departmentlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/department"
	"labyrinth/models/depposition"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (d DepartmentLogic) UpdateDepartment(
	userId,
	companyId uuid.UUID,
	updateDepartment *department.Department,
) error {
	// 1. Input validation
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "UpdateDepartment"),
			zap.Time("time", time.Now()),
		)
		return errors.New("user ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "UpdateDepartment"),
			zap.Time("time", time.Now()),
		)
		return errors.New("company ID cannot be empty")
	}

	if updateDepartment.ID == uuid.Nil {
		logger.NewWarnMessage("Empty department ID in update data",
			zap.String("operation", "UpdateDepartment"),
			zap.Time("time", time.Now()),
		)
		return errors.New("department ID cannot be empty")
	}

	if updateDepartment.CompanyID != companyId {
		logger.NewWarnMessage("Company ID mismatch",
			zap.String("operation", "UpdateDepartment"),
			zap.String("requested_company_id", companyId.String()),
			zap.String("department_company_id", updateDepartment.CompanyID.String()),
		)
		return errors.New("department does not belong to specified company")
	}

	if updateDepartment.ParentID == uuid.Nil {
		logger.NewWarnMessage("Empty parent department ID",
			zap.String("operation", "UpdateDepartment"),
			zap.Time("time", time.Now()),
		)
		return errors.New("parent department ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "UpdateDepartment"),
			zap.String("user_id", userId.String()),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin transaction (should be read-write for update operation)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "UpdateDepartment"),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// Ensure transaction is rolled back on error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "UpdateDepartment"),
				)
			}
		}
	}()

	ps := postgres.NewPostgresDB()

	// 5. Verify employee exists and belongs to company
	fetchedEmployee, err := ps.Employee.GetEmployeeByUserId(ctx, tx, userId, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Employee not found in company",
				zap.String("user_id", userId.String()),
				zap.String("company_id", companyId.String()),
			)
			return fmt.Errorf("employee not found in company: %w", err)
		}

		logger.NewErrMessage("Failed to fetch employee",
			zap.Error(err),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("failed to fetch employee: %w", err)
	}

	// 6. Get employee's main position
	fetchedPosition, err := ps.Position.GetPositionById(ctx, tx, fetchedEmployee.PositionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Position not found",
				zap.String("position_id", fetchedEmployee.PositionID.String()),
			)
			return fmt.Errorf("position not found: %w", err)
		}

		logger.NewErrMessage("Failed to fetch position",
			zap.Error(err),
			zap.String("position_id", fetchedEmployee.PositionID.String()),
		)
		return fmt.Errorf("failed to fetch position: %w", err)
	}

	// 7. Get employee's department relationship
	fetchedEmplDep, err := ps.DepartmentEmployee.GetEmployeeDepartmentByEmployeeId(ctx, tx, fetchedEmployee.ID, updateDepartment.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.NewErrMessage("Failed to fetch employee department relationship",
			zap.Error(err),
			zap.String("employee_id", fetchedEmployee.ID.String()),
		)
		return fmt.Errorf("failed to fetch employee department relationship: %w", err)
	}

	// 8. Check department position if employee belongs to department
	var fetchedEmpPosition *depposition.DepPosition
	if fetchedEmplDep != nil {
		fetchedEmpPosition, err = ps.DepartmentEmployeePosition.GetDepartmentPositionById(ctx, tx, fetchedEmplDep.PositionID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.NewErrMessage("Failed to fetch department position",
				zap.Error(err),
				zap.String("position_id", fetchedEmplDep.PositionID.String()),
			)
			return fmt.Errorf("failed to fetch department position: %w", err)
		}
	}

	// 9. Check update permissions
	hasPermission := (fetchedPosition != nil && fetchedPosition.Lvl <= 1) ||
		(fetchedEmpPosition != nil && fetchedEmpPosition.Level <= 1)

	if !hasPermission {
		logger.NewWarnMessage("Update permission denied",
			zap.String("user_id", userId.String()),
			zap.String("department_id", updateDepartment.ID.String()),
			zap.Int("company_position_level", fetchedPosition.Lvl),
			zap.Int("department_position_level", fetchedEmpPosition.Level),
		)
		return errors.New("insufficient permissions to update department")
	}

	// 10. Update department
	err = ps.Department.UpdateDepartment(ctx, tx, updateDepartment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Department not found",
				zap.String("department_id", updateDepartment.ID.String()),
			)
			return fmt.Errorf("department not found: %w", err)
		}

		logger.NewErrMessage("Failed to update department",
			zap.Error(err),
			zap.String("department_id", updateDepartment.ID.String()),
		)
		return fmt.Errorf("failed to update department: %w", err)
	}

	// 11. Commit transaction
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "UpdateDepartment"),
			zap.String("department_id", updateDepartment.ID.String()),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	// 12. Log successful update
	logger.NewInfoMessage("Department updated successfully",
		zap.String("department_id", updateDepartment.ID.String()),
		zap.String("updated_by", userId.String()),
		zap.Time("updated_at", time.Now()),
	)

	return nil
}
