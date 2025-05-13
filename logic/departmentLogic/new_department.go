package departmentlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/department"
	"labyrinth/models/depemployee"
	"labyrinth/models/depposition"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (d DepartmentLogic) NewDepartment(
	userId,
	companyId,
	parentId uuid.UUID,
	name,
	description string,
) (uuid.UUID, uuid.UUID, uuid.UUID, error) {
	// 1. Input validation
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "NewDepartment"),
			zap.Time("time", time.Now()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, errors.New("user ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "NewDepartment"),
			zap.Time("time", time.Now()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, errors.New("company ID cannot be empty")
	}

	if parentId == uuid.Nil {
		logger.NewWarnMessage("Empty parent department ID provided",
			zap.String("operation", "NewDepartment"),
			zap.Time("time", time.Now()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, errors.New("parent department ID cannot be empty")
	}

	checkName := strings.TrimSpace(name)
	checkDescription := strings.TrimSpace(description)
	if checkName == "" || checkDescription == "" {
		logger.NewWarnMessage("Empty name or description provided",
			zap.String("operation", "NewDepartment"),
			zap.Time("time", time.Now()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, errors.New("department name and description cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "NewDepartment"),
			zap.String("user_id", userId.String()),
			zap.String("company_id", companyId.String()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin transaction (should be read-write for create operations)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "NewDepartment"),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("transaction begin failed: %w", err)
	}

	// Ensure transaction is rolled back on error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "NewDepartment"),
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
			return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("employee not found in company: %w", err)
		}

		logger.NewErrMessage("Failed to fetch employee",
			zap.Error(err),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("failed to fetch employee: %w", err)
	}

	// 6. Generate UUIDs for new entities
	generatedDepId, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate department UUID",
			zap.Error(err),
			zap.String("operation", "NewDepartment"),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("failed to generate department UUID: %w", err)
	}

	generatedDepEmpId, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate department employee UUID",
			zap.Error(err),
			zap.String("operation", "NewDepartment"),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("failed to generate department employee UUID: %w", err)
	}

	generatedDepPosId, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate department position UUID",
			zap.Error(err),
			zap.String("operation", "NewDepartment"),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("failed to generate department position UUID: %w", err)
	}

	// 7. Create new department and related entities
	newDepartment := department.NewDepartment(generatedDepId, companyId, parentId, name, description)
	newDepEmployee := depemployee.NewDepartmentEmployee(generatedDepEmpId, fetchedEmployee.ID, generatedDepId, generatedDepPosId)
	newDepPosition := depposition.NewDepPosition(generatedDepPosId, generatedDepId, 0, "owner")

	// 8. Create department
	err = ps.Department.CreateDepartment(ctx, tx, newDepartment)
	if err != nil {
		logger.NewErrMessage("Failed to create department",
			zap.Error(err),
			zap.String("department_id", generatedDepId.String()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("failed to create department: %w", err)
	}

	// 9. Create department employee relationship
	err = ps.DepartmentEmployee.CreateEmployeeDepartment(ctx, tx, newDepEmployee)
	if err != nil {
		logger.NewErrMessage("Failed to create department employee",
			zap.Error(err),
			zap.String("department_id", generatedDepId.String()),
			zap.String("employee_id", fetchedEmployee.ID.String()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("failed to create department employee: %w", err)
	}

	// 10. Create department position
	err = ps.DepartmentEmployeePosition.CreateDepartmentPosition(ctx, tx, newDepPosition)
	if err != nil {
		logger.NewErrMessage("Failed to create department position",
			zap.Error(err),
			zap.String("department_id", generatedDepId.String()),
			zap.String("position_id", generatedDepPosId.String()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("failed to create department position: %w", err)
	}

	// 11. Commit transaction
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "NewDepartment"),
			zap.String("department_id", generatedDepId.String()),
		)
		return uuid.Nil, uuid.Nil, uuid.Nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	// 12. Log successful creation
	logger.NewInfoMessage("Department created successfully",
		zap.String("department_id", generatedDepId.String()),
		zap.String("company_id", companyId.String()),
		zap.String("created_by", userId.String()),
		zap.Time("created_at", time.Now()),
	)

	return generatedDepId, generatedDepEmpId, generatedDepPosId, nil
}
