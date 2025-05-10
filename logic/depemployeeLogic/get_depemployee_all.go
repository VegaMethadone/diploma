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

func (d DepemployeeLogic) GetAllDepEmployees(
	departmentId uuid.UUID,
) (*[]depemployee.DepartmentEmployee, error) {
	// 1. Validate input parameter
	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "GetAllDepEmployees"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("department ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "GetAllDepEmployees"),
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
			zap.String("operation", "GetAllDepEmployees"),
			zap.String("department_id", departmentId.String()),
		)
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}

	// Ensure proper transaction handling
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "GetAllDepEmployees"),
				)
			}
			return
		}
		// Commit for read-only transaction still needed to release resources
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "GetAllDepEmployees"),
			)
		}
	}()

	// 5. Fetch department employees
	ps := postgres.NewPostgresDB()
	fetchedDepEmplo, err := ps.DepartmentEmployee.GetEmployeesDepartmentByDepartmentId(ctx, tx, departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to get department employees",
			zap.Error(err),
			zap.String("operation", "GetAllDepEmployees"),
			zap.String("department_id", departmentId.String()),
		)
		return nil, fmt.Errorf("failed to get department employees: %w", err)
	}

	// 6. Check if any employees found
	if len(*fetchedDepEmplo) == 0 {
		logger.NewInfoMessage("No employees found for department",
			zap.String("operation", "GetAllDepEmployees"),
			zap.String("department_id", departmentId.String()),
		)
		return &[]depemployee.DepartmentEmployee{}, nil
	}

	logger.NewInfoMessage("Successfully retrieved department employees",
		zap.String("operation", "GetAllDepEmployees"),
		zap.String("department_id", departmentId.String()),
		zap.Int("employee_count", len(*fetchedDepEmplo)),
	)

	return fetchedDepEmplo, nil
}
