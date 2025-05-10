package departmentlogic

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

func (d DepartmentLogic) DeleteDepartment(userId, companyId, departmentId uuid.UUID) error {
	// 1. Input validation
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "DeleteDepartment"),
			zap.Time("time", time.Now()),
		)
		return errors.New("user ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "DeleteDepartment"),
			zap.Time("time", time.Now()),
		)
		return errors.New("company ID cannot be empty")
	}

	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "DeleteDepartment"),
			zap.Time("time", time.Now()),
		)
		return errors.New("department ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "DeleteDepartment"),
			zap.String("user_id", userId.String()),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin transaction (should be read-write for delete operation)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "DeleteDepartment"),
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
					zap.String("operation", "DeleteDepartment"),
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
	fetchedDepEmployee, err := ps.DepartmentEmployee.GetEmployeeDepartmentByEmployeeId(ctx, tx, fetchedEmployee.ID, departmentId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.NewErrMessage("Failed to fetch department employee relationship",
			zap.Error(err),
			zap.String("employee_id", fetchedEmployee.ID.String()),
		)
		return fmt.Errorf("failed to fetch department employee relationship: %w", err)
	}

	// 8. Check department position if employee belongs to department
	var fetchedDepPosition *depposition.DepPosition
	if fetchedDepEmployee != nil {
		fetchedDepPosition, err = ps.DepartmentEmployeePosition.GetDepartmentPositionById(ctx, tx, fetchedDepEmployee.PositionID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.NewErrMessage("Failed to fetch department position",
				zap.Error(err),
				zap.String("position_id", fetchedDepEmployee.PositionID.String()),
			)
			return fmt.Errorf("failed to fetch department position: %w", err)
		}
	}

	// 9. Check delete permissions (company admin or department admin)
	hasPermission := (fetchedPosition != nil && fetchedPosition.Lvl <= 1) ||
		(fetchedDepPosition != nil && fetchedDepPosition.Level <= 1)

	if !hasPermission {
		logger.NewWarnMessage("Delete permission denied",
			zap.String("user_id", userId.String()),
			zap.String("department_id", departmentId.String()),
			zap.Int("company_position_level", fetchedPosition.Lvl),
			zap.Int("department_position_level", fetchedDepPosition.Level),
		)
		return errors.New("insufficient permissions to delete department")
	}

	// 10. Verify department exists and belongs to company
	existingDep, err := ps.Department.GetDepartmentById(ctx, tx, departmentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Department not found",
				zap.String("department_id", departmentId.String()),
			)
			return fmt.Errorf("department not found: %w", err)
		}

		logger.NewErrMessage("Failed to fetch department",
			zap.Error(err),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("failed to fetch department: %w", err)
	}

	if existingDep.CompanyID != companyId {
		logger.NewWarnMessage("Department doesn't belong to company",
			zap.String("department_id", departmentId.String()),
			zap.String("department_company_id", existingDep.CompanyID.String()),
			zap.String("requested_company_id", companyId.String()),
		)
		return errors.New("department doesn't belong to specified company")
	}

	// 12. Delete department
	err = ps.Department.DeleteDepartment(ctx, tx, departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to delete department",
			zap.Error(err),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("failed to delete department: %w", err)
	}

	// 13. Commit transaction
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "DeleteDepartment"),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	// 14. Log successful deletion
	logger.NewInfoMessage("Department deleted successfully",
		zap.String("department_id", departmentId.String()),
		zap.String("deleted_by", userId.String()),
		zap.Time("deleted_at", time.Now()),
	)

	return nil
}
