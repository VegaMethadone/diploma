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

func (p PositionLogic) UpdatePosition(userId, companyId uuid.UUID, updatePosition *position.Position) error {
	// 1. Validate input parameters
	if userId == uuid.Nil {
		logger.NewWarnMessage("Empty user ID provided",
			zap.String("operation", "UpdatePosition"),
			zap.String("user_id", userId.String()),
		)
		return errors.New("user ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewWarnMessage("Empty company ID provided",
			zap.String("operation", "UpdatePosition"),
			zap.String("company_id", companyId.String()),
		)
		return errors.New("company ID cannot be empty")
	}

	if updatePosition == nil {
		logger.NewWarnMessage("Nil position provided",
			zap.String("operation", "UpdatePosition"),
			zap.String("user_id", userId.String()),
		)
		return errors.New("position data cannot be nil")
	}

	if updatePosition.ID == uuid.Nil {
		logger.NewWarnMessage("Empty position ID provided",
			zap.String("operation", "UpdatePosition"),
			zap.String("user_id", userId.String()),
		)
		return errors.New("position ID cannot be empty")
	}

	if strings.TrimSpace(updatePosition.Name) == "" {
		logger.NewWarnMessage("Empty position name provided",
			zap.String("operation", "UpdatePosition"),
			zap.String("position_id", updatePosition.ID.String()),
		)
		return errors.New("position name cannot be empty")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "UpdatePosition"),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("database connection failed: %w", err)
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
			zap.String("operation", "UpdatePosition"),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// 5. Deferred transaction handling
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "UpdatePosition"),
					zap.String("user_id", userId.String()),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "UpdatePosition"),
				zap.String("user_id", userId.String()),
			)
		}
	}()

	ps := postgres.NewPostgresDB()

	// 6. Verify employee access
	employee, err := ps.Employee.GetEmployeeByUserId(ctx, tx, userId, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Employee not found in company",
				zap.String("operation", "UpdatePosition"),
				zap.String("user_id", userId.String()),
				zap.String("company_id", companyId.String()),
			)
			return fmt.Errorf("employee not found in company: %w", err)
		}

		logger.NewErrMessage("Failed to get employee",
			zap.Error(err),
			zap.String("operation", "UpdatePosition"),
			zap.String("user_id", userId.String()),
		)
		return fmt.Errorf("failed to verify employee access: %w", err)
	}

	// 7. Check position level
	currentPosition, err := ps.Position.GetPositionById(ctx, tx, employee.PositionID)
	if err != nil {
		logger.NewErrMessage("Failed to get current position",
			zap.Error(err),
			zap.String("operation", "UpdatePosition"),
			zap.String("position_id", employee.PositionID.String()),
		)
		return fmt.Errorf("failed to get current position: %w", err)
	}

	if currentPosition.Lvl > 1 {
		logger.NewWarnMessage("Access denied - insufficient position level",
			zap.String("operation", "UpdatePosition"),
			zap.String("user_id", userId.String()),
			zap.Int("required_level", 1), // или 0
			zap.Int("current_level", currentPosition.Lvl),
		)
		return fmt.Errorf("access denied: position level %d is too low", currentPosition.Lvl)
	}

	// 8. Verify position belongs to company
	existingPosition, err := ps.Position.GetPositionById(ctx, tx, updatePosition.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Position not found",
				zap.String("operation", "UpdatePosition"),
				zap.String("position_id", updatePosition.ID.String()),
			)
			return fmt.Errorf("position not found: %w", err)
		}

		logger.NewErrMessage("Failed to get position",
			zap.Error(err),
			zap.String("operation", "UpdatePosition"),
			zap.String("position_id", updatePosition.ID.String()),
		)
		return fmt.Errorf("failed to verify position: %w", err)
	}

	if existingPosition.CompanyID != companyId {
		logger.NewWarnMessage("Position does not belong to company",
			zap.String("operation", "UpdatePosition"),
			zap.String("position_id", updatePosition.ID.String()),
			zap.String("company_id", companyId.String()),
		)
		return fmt.Errorf("position does not belong to company")
	}

	// 9. Update position
	err = ps.Position.UpdatePosition(ctx, tx, updatePosition)
	if err != nil {
		logger.NewErrMessage("Position update failed",
			zap.Error(err),
			zap.String("operation", "UpdatePosition"),
			zap.String("position_id", updatePosition.ID.String()),
		)
		return fmt.Errorf("failed to update position: %w", err)
	}

	// 10. Log success
	logger.NewInfoMessage("Position updated successfully",
		zap.String("operation", "UpdatePosition"),
		zap.String("user_id", userId.String()),
		zap.String("position_id", updatePosition.ID.String()),
		zap.String("new_name", updatePosition.Name),
		zap.Int("new_level", updatePosition.Lvl),
	)

	return nil
}
