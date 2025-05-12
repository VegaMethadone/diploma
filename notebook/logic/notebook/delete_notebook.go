package notebookLogic

import (
	"context"
	"fmt"
	m "labyrinth/database/mongo"
	"labyrinth/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (n NotebookMongoLogic) DeleteNotebook(notebookId uuid.UUID) error {
	// 1. Validate notebookId
	if notebookId == uuid.Nil {
		return fmt.Errorf("notebookId cannot be nil")
	}

	// 2. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Initialize MongoDB
	md, err := m.NewMongoDB()
	if err != nil {
		logger.NewErrMessage(
			"MongoDB initialization failed",
			zap.Error(err),
			zap.String("operation", "DeleteNotebook"),
			zap.String("notebook_id", notebookId.String()),
		)
		return fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	// 4. Start session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage(
			"MongoDB session start failed",
			zap.Error(err),
			zap.String("operation", "DeleteNotebook"),
			zap.String("notebook_id", notebookId.String()),
		)
		return fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	// 5. Delete notebook
	err = md.Notebook.DeleteNotebook(ctx, &session, notebookId.String())
	if err != nil {
		logger.NewErrMessage(
			"Failed to delete notebook",
			zap.Error(err),
			zap.String("operation", "DeleteNotebook"),
			zap.String("notebook_id", notebookId.String()),
		)
		return fmt.Errorf("failed to delete notebook: %w", err)
	}

	// 6. Delete related permissions
	err = md.Permission.DeletePermission(ctx, &session, notebookId.String())
	if err != nil {
		logger.NewErrMessage(
			"Failed to delete notebook permissions",
			zap.Error(err),
			zap.String("operation", "DeleteNotebook"),
			zap.String("notebook_id", notebookId.String()),
		)
		return fmt.Errorf("failed to delete notebook permissions: %w", err)
	}

	return nil
}
