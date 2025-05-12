package folderLogic

import (
	"context"
	"fmt"
	"labyrinth/database/mongo"
	"labyrinth/logger"
	"labyrinth/notebook/models/directory"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetFolder(folderId uuid.UUID) (*directory.Directory, error) {
	// 1. Validate input
	if folderId == uuid.Nil {
		logger.NewErrMessage("Invalid folder ID",
			zap.String("error", "folderId cannot be nil"),
			zap.String("operation", "GetFolder"),
		)
		return nil, fmt.Errorf("invalid folder ID: cannot be nil")
	}

	// 2. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Initialize MongoDB
	md, err := mongo.NewMongoDB()
	if err != nil {
		logger.NewErrMessage("MongoDB initialization failed",
			zap.Error(err),
			zap.String("operation", "GetFolder"),
			zap.String("folder_id", folderId.String()),
		)
		return nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	// 4. Start session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage("MongoDB session start failed",
			zap.Error(err),
			zap.String("operation", "GetFolder"),
			zap.String("folder_id", folderId.String()),
		)
		return nil, fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	// 5. Fetch folder
	fetchedDir, err := md.Folder.GetFolderByFolderId(ctx, &session, folderId.String())
	if err != nil {
		logger.NewErrMessage("Failed to get folder",
			zap.Error(err),
			zap.String("operation", "GetFolder"),
			zap.String("folder_id", folderId.String()),
		)
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	// 6. Log success
	logger.NewInfoMessage("Folder retrieved successfully",
		zap.String("folder_id", folderId.String()),
		zap.String("operation", "GetFolder"),
	)

	return fetchedDir, nil
}
