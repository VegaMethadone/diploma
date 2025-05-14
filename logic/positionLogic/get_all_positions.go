package positionlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/position"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (p PositionLogic) GetAllPositions(userId, companyId uuid.UUID) (*[]position.Position, error) {
	// 1. Validate input parameters
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "GetAllPositions"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("user ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "GetAllPositions"),
			zap.Time("time", time.Now()),
		)
		return nil, errors.New("company ID cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "GetAllPositions"),
			zap.String("user_id", userId.String()),
		)
		return nil, fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin transaction (should be read-only for query operation)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "GetAllPositions"),
			zap.String("user_id", userId.String()),
		)
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}

	// 5. Deferred transaction handling
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "GetAllPositions"),
					zap.String("user_id", userId.String()),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "GetAllPositions"),
				zap.String("user_id", userId.String()),
			)
		}
	}()

	ps := postgres.NewPostgresDB()

	// 6. Check employee access
	_, err = ps.Employee.GetEmployeeByUserId(ctx, tx, userId, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Employee not found in company",
				zap.String("operation", "GetAllPositions"),
				zap.String("user_id", userId.String()),
				zap.String("company_id", companyId.String()),
			)
			return nil, fmt.Errorf("employee not found in company: %w", err)
		}

		logger.NewErrMessage("Failed to get employee",
			zap.Error(err),
			zap.String("operation", "GetAllPositions"),
			zap.String("user_id", userId.String()),
		)
		return nil, fmt.Errorf("failed to verify employee access: %w", err)
	}

	// 7. Get all positions for company
	positions, err := ps.Position.GetPositionsByCompanyId(ctx, tx, companyId)
	if err != nil {
		logger.NewErrMessage("Failed to get positions",
			zap.Error(err),
			zap.String("operation", "GetAllPositions"),
			zap.String("company_id", companyId.String()),
		)
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// 8. Log success
	logger.NewInfoMessage("Positions retrieved successfully",
		zap.String("operation", "GetAllPositions"),
		zap.String("user_id", userId.String()),
		zap.String("company_id", companyId.String()),
		zap.Int("positions_count", len(*positions)),
	)

	return positions, nil
}
