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

func (e EmployeeLogic) GetAllEmployee(
	companyId uuid.UUID,
) (*[]employee.Employee, error) {
	// 1. Validate input
	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "GetAllEmployee"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("company ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "GetAllEmployee"),
			zap.String("company_id", companyId.String()),
		)
		return nil, fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin transaction (read-only for fetching data)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "GetAllEmployee"),
			zap.String("company_id", companyId.String()),
		)
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	// 5. Fetch employees from database
	ps := postgres.NewPostgresDB()
	fetchedEmployees, err := ps.Employee.GetEmployeesByCompanyId(ctx, tx, companyId)
	if err != nil {
		logger.NewErrMessage("Failed to get employees",
			zap.Error(err),
			zap.String("operation", "GetAllEmployee"),
			zap.String("company_id", companyId.String()),
		)
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	logger.NewInfoMessage("Successfully retrieved employees",
		zap.Int("count", len(*fetchedEmployees)),
		zap.String("operation", "GetAllEmployee"),
		zap.String("company_id", companyId.String()),
	)

	return fetchedEmployees, nil
}
