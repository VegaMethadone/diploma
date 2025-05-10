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

func (d DepemploeePosLogic) GetAllDepEmployeePos(
	departmentId uuid.UUID,
) (*[]depposition.DepPosition, error) {
	// 1. Validate input parameter
	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "GetAllDepEmployeePos"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("department ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "GetAllDepEmployeePos"),
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
			zap.String("operation", "GetAllDepEmployeePos"),
			zap.String("department_id", departmentId.String()),
		)
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}

	// 5. Ensure proper transaction handling
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "GetAllDepEmployeePos"),
				)
			}
			return
		}
	}()

	// 6. Fetch department positions
	ps := postgres.NewPostgresDB()
	fetchedDepPos, err := ps.DepartmentEmployeePosition.GetDepartmentPositionsByDepartmentId(ctx, tx, departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to get department positions",
			zap.Error(err),
			zap.String("operation", "GetAllDepEmployeePos"),
			zap.String("department_id", departmentId.String()),
		)
		return nil, fmt.Errorf("failed to get department positions: %w", err)
	}

	logger.NewInfoMessage("Successfully retrieved department positions",
		zap.String("operation", "GetAllDepEmployeePos"),
		zap.String("department_id", departmentId.String()),
		zap.Int("positions_count", len(*fetchedDepPos)),
	)

	return fetchedDepPos, nil
}
