package folderLogic

import (
	"context"
	"database/sql"
	"fmt"
	m "labyrinth/database/mongo"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/notebook/models/directory"
	"labyrinth/notebook/models/permission"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (f FolderMongoLogic) CreateFolder(employeeId, companyId, divisionId, parentId uuid.UUID, isPrimary bool, title, description string) error {
	// 1. Validate input parameters
	if employeeId == uuid.Nil {
		return fmt.Errorf("employeeId cannot be nil")
	}
	if companyId == uuid.Nil {
		return fmt.Errorf("companyId cannot be nil")
	}
	if divisionId == uuid.Nil {
		return fmt.Errorf("divisionId cannot be nil")
	}
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("title cannot be empty")
	}

	// 2. Initialize PostgreSQL connection
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "CreateFolder"),
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
			zap.String("operation", "CreateFolder"),
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
					zap.String("operation", "CreateFolder"),
				)
			}
		}
	}()

	// 6. Initialize MongoDB
	md, err := m.NewMongoDB()
	if err != nil {
		logger.NewErrMessage("MongoDB initialization failed",
			zap.Error(err),
			zap.String("operation", "CreateFolder"),
		)
		return fmt.Errorf("mongodb initialization failed: %w", err)
	}

	// 7. Generate and validate UUID
	generatedId, err := postgres.NewPostgresDB().UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("UUID generation failed",
			zap.Error(err),
			zap.String("operation", "CreateFolder"),
		)
		return fmt.Errorf("uuid generation failed: %w", err)
	}

	// 8. Start MongoDB session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage("MongoDB session start failed",
			zap.Error(err),
			zap.String("operation", "CreateFolder"),
		)
		return fmt.Errorf("mongodb session start failed: %w", err)
	}
	defer session.EndSession(ctx)

	// 9. Create new folder
	newFolder := directory.NewDirectory(employeeId, companyId, divisionId, generatedId, parentId, "1.0.0", isPrimary, title, description)
	err = md.Folder.CreateFolder(ctx, &session, &newFolder)
	if err != nil {
		logger.NewErrMessage("Folder creation failed",
			zap.Error(err),
			zap.String("operation", "CreateFolder"),
			zap.String("folder_id", generatedId.String()),
		)
		return fmt.Errorf("folder creation failed: %w", err)
	}

	newPermission := permission.NewPermission(employeeId.String(), generatedId.String(), generatedId.String(), "folder", newFolder.ID)
	err = md.Permission.CreatePermission(ctx, &session, &newPermission)
	if err != nil {
		logger.NewErrMessage("Permission creation failed",
			zap.Error(err),
			zap.String("operation", "CreateFolder"),
			zap.String("folder_id", generatedId.String()),
		)
		return fmt.Errorf("permission creation failed: %w", err)
	}

	// 10. Commit transaction if everything succeeded
	if err = tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "CreateFolder"),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
