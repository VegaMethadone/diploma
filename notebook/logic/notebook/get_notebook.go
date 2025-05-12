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

func (n NotebookMongoLogic) GetNotebook(notebookId uuid.UUID) (*journal.Notebook, error) {
	if notebookId == uuid.Nil {
		logger.NewErrMessage("Empty notebook ID provided",
			zap.String("operation", "GetNotebook"),
		)
		return nil, errors.New("notebook ID cannot be empty")
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
		return nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	// 4. Start session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage("MongoDB session start failed",
			zap.Error(err),
			zap.String("operation", "GetNotebook"),
			zap.String("notebook_id", notebookId.String()),
		)
		return nil, fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	fetchedNotebook, err := md.Notebook.GetNotebookById(ctx, &session, notebookId.String())
	if err != nil {
		logger.NewErrMessage("Failed to get notebook",
			zap.Error(err),
			zap.String("operation", "GetNotebook"),
			zap.String("notebook_id", notebookId.String()),
		)
		return nil, fmt.Errorf("failed to create notebook: %w", err)
	}

	logger.NewInfoMessage("Notebook get successfully",
		zap.String("operation", "GetNotebook"),
		zap.String("notebook_id", notebookId.String()),
	)

	return fetchedNotebook, nil
}
