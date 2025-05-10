package depemployeeposlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/depposition"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (p DepemploeePosLogic) NewDepemployeePos(
	departmentId uuid.UUID,
	lvl int,
	name string,
) error {
	// 1. Validate input parameters
	if departmentId == uuid.Nil {
		logger.NewWarnMessage("Empty department ID provided",
			zap.String("operation", "NewDepemployeePos"),
			zap.Time("time", time.Now()),
		)
		return errors.New("department ID cannot be empty")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		logger.NewWarnMessage("Empty position name provided",
			zap.String("operation", "NewDepemployeePos"),
		)
		return errors.New("position name cannot be empty")
	}

	if lvl <= 0 {
		logger.NewWarnMessage("Invalid position level provided",
			zap.String("operation", "NewDepemployeePos"),
			zap.Int("level", lvl),
		)
		return errors.New("position level must be positive")
	}

	// 2. Initialize database connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "NewDepemployeePos"),
			zap.String("department_id", departmentId.String()),
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
			zap.String("operation", "NewDepemployeePos"),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// 5. Ensure proper transaction handling
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "NewDepemployeePos"),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "NewDepemployeePos"),
			)
		}
	}()

	// 6. Generate and validate new UUID
	ps := postgres.NewPostgresDB()
	generatedId, err := ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("Failed to generate UUID",
			zap.Error(err),
			zap.String("operation", "NewDepemployeePos"),
			zap.String("department_id", departmentId.String()),
		)
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	// 7. Create new department position
	newDepPos := depposition.NewDepPosition(generatedId, departmentId, lvl, name)

	// 8. Save to database
	err = ps.DepartmentEmployeePosition.CreateDepartmentPosition(ctx, tx, newDepPos)
	if err != nil {
		logger.NewErrMessage("Failed to create department position",
			zap.Error(err),
			zap.String("operation", "NewDepemployeePos"),
			zap.String("department_id", departmentId.String()),
			zap.String("position_id", generatedId.String()),
			zap.String("position_name", name),
			zap.Int("position_level", lvl),
		)
		return fmt.Errorf("failed to create department position: %w", err)
	}

	logger.NewInfoMessage("Successfully created new department position",
		zap.String("operation", "NewDepemployeePos"),
		zap.String("department_id", departmentId.String()),
		zap.String("position_id", generatedId.String()),
		zap.String("position_name", name),
		zap.Int("position_level", lvl),
	)

	return nil
}
