package employeelogic

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

func (e EmployeeLogic) DeleteEmployee(
	userId,
	companyId,
	employeeId uuid.UUID,
) error {
	// 1. Input validation
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "DeleteEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("user ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "DeleteEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("company ID cannot be empty")
	}

	if employeeId == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID provided",
			zap.String("operation", "DeleteEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "DeleteEmployee"),
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
			zap.String("operation", "DeleteEmployee"),
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
					zap.String("operation", "DeleteEmployee"),
				)
			}
		}
	}()

	ps := postgres.NewPostgresDB()

	// 5. Verify requesting employee exists and belongs to company
	fetchedEmployee, err := ps.Employee.GetEmployeeByUserId(ctx, tx, userId, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Requesting employee not found in company",
				zap.String("user_id", userId.String()),
				zap.String("company_id", companyId.String()),
			)
			return fmt.Errorf("requesting employee not found in company: %w", err)
		}

		logger.NewErrMessage("Failed to fetch requesting employee",
			zap.Error(err),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("failed to fetch requesting employee: %w", err)
	}

	// 6. Prevent self-deletion
	if fetchedEmployee.ID == employeeId {
		logger.NewWarnMessage("Attempt to delete self",
			zap.String("user_id", userId.String()),
			zap.String("employee_id", employeeId.String()),
		)
		return errors.New("cannot delete yourself")
	}

	// 7. Get requesting employee's position
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

	// 8. Check if requesting employee has sufficient privileges (level <= 1)
	if fetchedPosition.Lvl > 1 {
		logger.NewWarnMessage("Insufficient privileges to delete employee",
			zap.String("user_id", userId.String()),
			zap.Int("position_level", fetchedPosition.Lvl),
		)
		return errors.New("insufficient privileges to delete employee")
	}

	// 9. Verify target employee exists and belongs to company
	targetEmployee, err := ps.Employee.GetEmployeeByUserId(ctx, tx, userId, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Target employee not found",
				zap.String("employee_id", employeeId.String()),
			)
			return fmt.Errorf("target employee not found: %w", err)
		}

		logger.NewErrMessage("Failed to fetch target employee",
			zap.Error(err),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("failed to fetch target employee: %w", err)
	}

	if targetEmployee.CompanyID != companyId {
		logger.NewWarnMessage("Target employee doesn't belong to company",
			zap.String("employee_id", employeeId.String()),
			zap.String("employee_company_id", targetEmployee.CompanyID.String()),
			zap.String("requested_company_id", companyId.String()),
		)
		return errors.New("target employee doesn't belong to specified company")
	}

	// 10. Delete employee
	err = ps.Employee.DeleteEmployee(ctx, tx, employeeId)
	if err != nil {
		logger.NewErrMessage("Failed to delete employee",
			zap.Error(err),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	// 11. Commit transaction
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "DeleteEmployee"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	// 12. Log successful deletion
	logger.NewInfoMessage("Employee deleted successfully",
		zap.String("employee_id", employeeId.String()),
		zap.String("deleted_by", userId.String()),
		zap.Time("deleted_at", time.Now()),
	)

	return nil
}
