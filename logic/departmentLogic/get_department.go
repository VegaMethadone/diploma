package departmentlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/department"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (d DepartmentLogic) GetDepartment(userId, companyId, departmentId uuid.UUID) (*department.Department, error) {
	// 1. Input validation
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "GetDepartment"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("user ID cannot be empty")
	}

	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "GetDepartment"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("department ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "GetDepartment"),
			zap.String("user_id", userId.String()),
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
			zap.String("operation", "GetDepartment"),
			zap.String("user_id", userId.String()),
		)
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}
	defer tx.Rollback() // Safe rollback for read-only transaction

	ps := postgres.NewPostgresDB()

	// 5. Verify employee exists and belongs to company
	fetchedEmployee, err := ps.Employee.GetEmployeeByUserId(ctx, tx, userId, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Employee not found in company",
				zap.String("user_id", userId.String()),
				zap.String("company_id", companyId.String()),
			)
			return nil, fmt.Errorf("employee not found in company: %w", err)
		}

		logger.NewErrMessage("Failed to fetch employee",
			zap.Error(err),
			zap.String("user_id", userId.String()),
		)
		return nil, fmt.Errorf("failed to fetch employee: %w", err)
	}

	// 6. Get employee position
	fetchedPosition, err := ps.Position.GetPositionById(ctx, tx, fetchedEmployee.PositionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Position not found",
				zap.String("position_id", fetchedEmployee.PositionID.String()),
			)
			return nil, fmt.Errorf("position not found: %w", err)
		}

		logger.NewErrMessage("Failed to fetch position",
			zap.Error(err),
			zap.String("position_id", fetchedEmployee.PositionID.String()),
		)
		return nil, fmt.Errorf("failed to fetch position: %w", err)
	}

	// 7. Check if employee belongs to department
	employeeDepExists, err := ps.DepartmentEmployee.ExistsEmployeeDepartment(ctx, tx, fetchedEmployee.ID, departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to check department employee relationship",
			zap.Error(err),
			zap.String("employee_id", fetchedEmployee.ID.String()),
			zap.String("department_id", departmentId.String()),
		)
		return nil, fmt.Errorf("failed to check department employee relationship: %w", err)
	}

	// 8. Check access rights (position level <= 1 or employee belongs to department)
	if fetchedPosition.Lvl <= 1 || employeeDepExists {
		fetchedDep, err := ps.Department.GetDepartmentById(ctx, tx, departmentId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.NewWarnMessage("Department not found",
					zap.String("department_id", departmentId.String()),
				)
				return nil, fmt.Errorf("department not found: %w", err)
			}

			logger.NewErrMessage("Failed to fetch department",
				zap.Error(err),
				zap.String("department_id", departmentId.String()),
			)
			return nil, fmt.Errorf("failed to fetch department: %w", err)
		}

		// 9. Verify department belongs to company
		if fetchedDep.CompanyID != companyId {
			logger.NewWarnMessage("Department doesn't belong to company",
				zap.String("department_id", departmentId.String()),
				zap.String("department_company_id", fetchedDep.CompanyID.String()),
				zap.String("requested_company_id", companyId.String()),
			)
			return nil, errors.New("department doesn't belong to specified company")
		}

		// 10. Log successful access
		logger.NewInfoMessage("Department accessed successfully",
			zap.String("department_id", departmentId.String()),
			zap.String("accessed_by", userId.String()),
			zap.Time("access_time", time.Now()),
		)

		return fetchedDep, nil
	}

	// 11. Access denied
	logger.NewWarnMessage("Access to department denied",
		zap.String("user_id", userId.String()),
		zap.String("department_id", departmentId.String()),
		zap.Int("position_level", fetchedPosition.Lvl),
		zap.Bool("department_employee", employeeDepExists),
	)
	return nil, errors.New("access to department denied")
}
