package positionlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/position"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (p PositionLogic) NewPosition(userId, companyId uuid.UUID, lvl int, name string) (uuid.UUID, error) {
	// 1. Validate input parameters
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "NewPosition"),
			zap.Time("time", time.Now()),
		)
		return uuid.Nil, errors.New("user ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "NewPosition"),
			zap.Time("time", time.Now()),
		)
		return uuid.Nil, errors.New("company ID cannot be empty")
	}

	if strings.TrimSpace(name) == "" {
		logger.NewWarnMessage("Empty position name provided",
			zap.String("operation", "NewPosition"),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, errors.New("position name cannot be empty")
	}

	if lvl < 0 {
		logger.NewWarnMessage("Invalid position level provided",
			zap.String("operation", "NewPosition"),
			zap.Int("level", lvl),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, errors.New("position level must be positive")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "NewPosition"),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin transaction
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "NewPosition"),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, fmt.Errorf("transaction begin failed: %w", err)
	}

	// 5. Deferred transaction handling
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "NewPosition"),
					zap.String("user_id", userId.String()),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "NewPosition"),
				zap.String("user_id", userId.String()),
			)
		}
	}()

	// 6. Generate and reserve UUID
	ps := postgres.NewPostgresDB()
	generatedId, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("UUID generation failed",
			zap.Error(err),
			zap.String("operation", "NewPosition"),
			zap.String("user_id", userId.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	// 7. Create new position
	newPosition := position.NewPosition(generatedId, companyId, lvl, name)
	err = ps.Position.CreatePosition(ctx, tx, &newPosition)
	if err != nil {
		logger.NewErrMessage("Position creation failed",
			zap.Error(err),
			zap.String("operation", "NewPosition"),
			zap.String("user_id", userId.String()),
			zap.String("position_id", generatedId.String()),
		)
		return uuid.Nil, fmt.Errorf("failed to create position: %w", err)
	}

	// 8. Log success
	logger.NewInfoMessage("Position created successfully",
		zap.String("operation", "NewPosition"),
		zap.String("user_id", userId.String()),
		zap.String("company_id", companyId.String()),
		zap.String("position_id", generatedId.String()),
		zap.String("position_name", name),
		zap.Int("position_level", lvl),
	)

	return generatedId, nil
}
