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

func (e EmployeeLogic) UpdateEmployee(
	userId,
	companyId uuid.UUID,
	updatedEmployee *employee.Employee,
) error {
	// 1. Input validation
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "UpdateEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("user ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "UpdateEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("company ID cannot be empty")
	}

	if updatedEmployee.ID == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID in request",
			zap.String("operation", "UpdateEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee ID cannot be empty")
	}

	if updatedEmployee.PositionID == uuid.Nil {
		logger.NewWarnMessage("Empty position ID in request",
			zap.String("operation", "UpdateEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("position ID cannot be empty")
	}

	if updatedEmployee.CompanyID != companyId {
		logger.NewWarnMessage("Company ID mismatch",
			zap.String("operation", "UpdateEmployee"),
			zap.String("requested_company_id", companyId.String()),
			zap.String("employee_company_id", updatedEmployee.CompanyID.String()),
		)
		return errors.New("employee does not belong to specified company")
	}

	if updatedEmployee.UserID == uuid.Nil {
		logger.NewWarnMessage("Empty user ID in employee data",
			zap.String("operation", "UpdateEmployee"),
			zap.Time("time", time.Now()),
		)
		return errors.New("employee user ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "UpdateEmployee"),
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
			zap.String("operation", "UpdateEmployee"),
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
					zap.String("operation", "UpdateEmployee"),
				)
			}
		}
	}()

	ps := postgres.NewPostgresDB()

	// 5. Verify position belongs to company
	fetchedPosition, err := ps.Position.GetPositionById(ctx, tx, updatedEmployee.PositionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Position not found",
				zap.String("position_id", updatedEmployee.PositionID.String()),
			)
			return fmt.Errorf("position not found: %w", err)
		}

		logger.NewErrMessage("Failed to fetch position",
			zap.Error(err),
			zap.String("position_id", updatedEmployee.PositionID.String()),
		)
		return fmt.Errorf("failed to fetch position: %w", err)
	}

	if fetchedPosition.CompanyID != companyId {
		logger.NewWarnMessage("Position doesn't belong to company",
			zap.String("position_id", updatedEmployee.PositionID.String()),
			zap.String("position_company_id", fetchedPosition.CompanyID.String()),
			zap.String("requested_company_id", companyId.String()),
		)
		return errors.New("position doesn't belong to specified company")
	}

	// 6. Verify employee exists and belongs to company
	fetchedEmployee, err := ps.Employee.GetEmployeeByUserId(ctx, tx, userId, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Employee not found",
				zap.String("user_id", userId.String()),
			)
			return fmt.Errorf("employee not found: %w", err)
		}

		logger.NewErrMessage("Failed to fetch employee",
			zap.Error(err),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("failed to fetch employee: %w", err)
	}

	if fetchedEmployee.CompanyID != companyId {
		logger.NewWarnMessage("Employee doesn't belong to company",
			zap.String("employee_id", fetchedEmployee.ID.String()),
			zap.String("employee_company_id", fetchedEmployee.CompanyID.String()),
			zap.String("requested_company_id", companyId.String()),
		)
		return errors.New("employee doesn't belong to specified company")
	}

	// 7. Update employee data
	err = ps.Employee.UpdateEmployee(ctx, tx, updatedEmployee)
	if err != nil {
		logger.NewErrMessage("Failed to update employee",
			zap.Error(err),
			zap.String("employee_id", updatedEmployee.ID.String()),
		)
		return fmt.Errorf("failed to update employee: %w", err)
	}

	// 8. Commit transaction
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "UpdateEmployee"),
			zap.String("employee_id", updatedEmployee.ID.String()),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	// 9. Log successful update
	logger.NewInfoMessage("Employee updated successfully",
		zap.String("employee_id", updatedEmployee.ID.String()),
		zap.String("user_id", userId.String()),
		zap.String("company_id", companyId.String()),
		zap.Time("updated_at", time.Now()),
	)

	return nil
}
