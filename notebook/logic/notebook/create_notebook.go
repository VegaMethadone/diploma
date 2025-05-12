package notebookLogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/mongo"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/notebook/models/journal"
	"labyrinth/notebook/models/permission"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (n NotebookMongoLogic) NewNotebook(employeeId, companyId, divisionId uuid.UUID, title, description string) error {
	// 1. Validate input parameters
	if employeeId == uuid.Nil {
		logger.NewErrMessage("Empty employee ID provided",
			zap.String("operation", "NewNotebook"),
		)
		return errors.New("employee ID cannot be empty")
	}

	if companyId == uuid.Nil {
		logger.NewErrMessage("Empty company ID provided",
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
		)
		return errors.New("company ID cannot be empty")
	}

	if divisionId == uuid.Nil {
		logger.NewErrMessage("Empty division ID provided",
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
		)
		return errors.New("division ID cannot be empty")
	}

	if strings.TrimSpace(title) == "" {
		logger.NewErrMessage("Empty notebook title provided",
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
		)
		return errors.New("notebook title cannot be empty")
	}

	// 2. Initialize PostgreSQL connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Begin transaction
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	// 5. Transaction rollback handler
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.NewErrMessage("Transaction rollback failed",
					zap.Error(rbErr),
					zap.String("operation", "NewNotebook"),
					zap.String("employee_id", employeeId.String()),
				)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			logger.NewErrMessage("Transaction commit failed",
				zap.Error(err),
				zap.String("operation", "NewNotebook"),
				zap.String("employee_id", employeeId.String()),
			)
		}
	}()

	// 6. Generate and validate UUID
	generatedId, err := postgres.NewPostgresDB().UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("UUID generation failed",
			zap.Error(err),
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("uuid generation failed: %w", err)
	}

	// 7. Initialize MongoDB
	md, err := mongo.NewMongoDB()
	if err != nil {
		logger.NewErrMessage("MongoDB initialization failed",
			zap.Error(err),
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	// 8. Start MongoDB session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage("MongoDB session start failed",
			zap.Error(err),
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
		)
		return fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	// 9. Create new notebook
	newNotebook := journal.NewNotebook(
		employeeId.String(),
		companyId.String(),
		divisionId.String(),
		generatedId.String(),
		title,
		description,
	)

	err = md.Notebook.CreateNotebook(ctx, &session, &newNotebook)
	if err != nil {
		logger.NewErrMessage("Failed to create notebook",
			zap.Error(err),
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
			zap.String("notebook_id", generatedId.String()),
		)
		return fmt.Errorf("failed to create notebook: %w", err)
	}

	// 10. Create permission for the notebook
	newPerm := permission.NewPermission(
		employeeId.String(),
		generatedId.String(),
		generatedId.String(),
		"file",
		newNotebook.ID,
	)

	err = md.Permission.CreatePermission(ctx, &session, &newPerm)
	if err != nil {
		logger.NewErrMessage("Failed to create permission",
			zap.Error(err),
			zap.String("operation", "NewNotebook"),
			zap.String("employee_id", employeeId.String()),
			zap.String("notebook_id", generatedId.String()),
		)
		return fmt.Errorf("failed to create permission: %w", err)
	}

	logger.NewInfoMessage("Notebook created successfully",
		zap.String("operation", "NewNotebook"),
		zap.String("employee_id", employeeId.String()),
		zap.String("notebook_id", generatedId.String()),
		zap.String("company_id", companyId.String()),
		zap.String("division_id", divisionId.String()),
	)

	return nil
}
