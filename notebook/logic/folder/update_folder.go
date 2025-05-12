package folderLogic

import (
	"context"
	"fmt"
	m "labyrinth/database/mongo"
	"labyrinth/logger"
	"labyrinth/notebook/models/directory"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (f FolderMongoLogic) UpdateFolder(folderId uuid.UUID, dir *directory.Directory) error {
	// 1. Validate input parameters
	if folderId == uuid.Nil {
		return fmt.Errorf("folderId cannot be nil")
	}
	if dir == nil {
		return fmt.Errorf("directory object cannot be nil")
	}

	// 2. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Initialize MongoDB
	md, err := m.NewMongoDB()
	if err != nil {
		logger.NewErrMessage("MongoDB initialization failed",
			zap.Error(err),
			zap.String("operation", "UpdateFolder"),
			zap.String("folder_id", folderId.String()),
		)
		return fmt.Errorf("mongodb initialization failed: %w", err)
	}

	// 4. Start MongoDB session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage("MongoDB session start failed",
			zap.Error(err),
			zap.String("operation", "UpdateFolder"),
			zap.String("folder_id", folderId.String()),
		)
		return fmt.Errorf("mongodb session start failed: %w", err)
	}
	defer session.EndSession(ctx)

	// 5. Update folder in MongoDB
	err = md.Folder.UpdateFolder(ctx, &session, folderId.String(), dir)
	if err != nil {
		logger.NewErrMessage("Folder update failed",
			zap.Error(err),
			zap.String("operation", "UpdateFolder"),
			zap.String("folder_id", folderId.String()),
		)
		return fmt.Errorf("failed to update folder: %w", err)
	}

	return nil
}
