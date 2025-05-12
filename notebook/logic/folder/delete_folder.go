package folderLogic

import (
	"context"
	"fmt"
	m "labyrinth/database/mongo"
	"labyrinth/logger"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (f FolderMongoLogic) DeleteFolder(folderId uuid.UUID) error {
	// 1. Validate input
	if folderId == uuid.Nil {
		return fmt.Errorf("folderId cannot be nil")
	}

	// 2. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Initialize MongoDB
	md, err := m.NewMongoDB()
	if err != nil {
		logger.NewErrMessage("MongoDB initialization failed",
			zap.Error(err),
			zap.String("operation", "DeleteFolder"),
			zap.String("folder_id", folderId.String()),
		)
		return fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	// 4. Start session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage("MongoDB session start failed",
			zap.Error(err),
			zap.String("operation", "DeleteFolder"),
			zap.String("folder_id", folderId.String()),
		)
		return fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	// 5. Execute recursive deletion
	if err := cleanDir(ctx, &session, md, folderId.String()); err != nil {
		logger.NewErrMessage("Failed to delete folder",
			zap.Error(err),
			zap.String("operation", "DeleteFolder"),
			zap.String("folder_id", folderId.String()),
		)
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	logger.NewInfoMessage("Folder deleted successfully",
		zap.String("folder_id", folderId.String()),
		zap.String("operation", "DeleteFolder"),
	)

	return nil
}

func cleanDir(ctx context.Context, session *mongo.Session, md *m.MongoDB, dirId string) error {
	dir, err := md.Folder.GetFolderByFolderId(ctx, session, dirId)
	if err != nil {
		return fmt.Errorf("failed to get directory %s: %w", dirId, err)
	}

	for _, folder := range dir.Folders {
		if err := cleanDir(ctx, session, md, folder.FolderUUID); err != nil {
			return fmt.Errorf("failed to clean subdirectory %s: %w", folder.FolderUUID, err)
		}
	}

	for _, file := range dir.Files {
		if err := md.Notebook.DeleteNotebook(ctx, session, file.FileUUID); err != nil {
			return fmt.Errorf("failed to delete notebook %s: %w", file.FileUUID, err)
		}
		if err := md.Permission.DeletePermission(ctx, session, file.FileUUID); err != nil {
			return fmt.Errorf("failed to delete notebook permission %s: %w", file.FileUUID, err)
		}
	}

	if err := md.Folder.DeleteFolder(ctx, session, dirId); err != nil {
		return fmt.Errorf("failed to delete directory %s: %w", dirId, err)
	}

	if err := md.Permission.DeletePermission(ctx, session, dir.UuidID); err != nil {
		return fmt.Errorf("failed to delete directory permission %s: %w", dirId, err)
	}

	return nil
}
