package notebookLogic

import (
	"context"
	"errors"
	"fmt"
	m "labyrinth/database/mongo"
	"labyrinth/logger"
	"labyrinth/notebook/models/journal"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (n NotebookMongoLogic) UpdateNotebook(notebookId uuid.UUID, updatedNotebook *journal.Notebook) error {
	if notebookId == uuid.Nil {
		logger.NewErrMessage("Empty notebook ID provided",
			zap.String("operation", "UpdateNotebook"),
		)
		return errors.New("notebook ID cannot be empty")
	}
	if updatedNotebook == nil {
		logger.NewErrMessage("Empty notebook provided",
			zap.String("operation", "UpdateNotebook"),
		)
		return errors.New("notebook cannot be empty")
	}

	// 2. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Initialize MongoDB
	md, err := m.NewMongoDB()
	if err != nil {
		logger.NewErrMessage("MongoDB initialization failed",
			zap.Error(err),
			zap.String("operation", "GetNotebook"),
			zap.String("notebook_id", notebookId.String()),
		)
		return fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	// 4. Start session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage("MongoDB session start failed",
			zap.Error(err),
			zap.String("operation", "GetNotebook"),
			zap.String("notebook_id", notebookId.String()),
		)
		return fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	err = md.Notebook.UpdateNotebook(ctx, &session, notebookId.String(), updatedNotebook)
	if err != nil {
		logger.NewErrMessage("update notebook failed",
			zap.Error(err),
			zap.String("operation", "UpdateNotebook"),
			zap.String("notebook_id", notebookId.String()),
		)
		return fmt.Errorf("failed to update notebook: %w", err)
	}

	logger.NewInfoMessage("Notebook created successfully",
		zap.String("operation", "GetNotebook"),
		zap.String("notebook_id", notebookId.String()),
	)

	return nil
}
