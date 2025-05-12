package permissionLogic

import (
	"context"
	"errors"
	"fmt"
	"labyrinth/database/mongo"
	"labyrinth/logger"
	"labyrinth/notebook/models/permission"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (p PermissionMongoLogic) UpdatePermission(objectId uuid.UUID, updatedPerm *permission.Permission) error {
	// 1. Validate input parameters
	if objectId == uuid.Nil {
		logger.NewErrMessage("Invalid permission ID",
			zap.String("operation", "UpdatePermission"),
			zap.Error(errors.New("permission ID cannot be nil")),
		)
		return errors.New("invalid permission ID: cannot be nil")
	}

	if updatedPerm == nil {
		logger.NewErrMessage("Nil permission provided",
			zap.String("operation", "UpdatePermission"),
			zap.String("permission_id", objectId.String()),
		)
		return errors.New("permission data cannot be nil")
	}

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Initialize MongoDB
	md, err := mongo.NewMongoDB()
	if err != nil {
		logger.NewErrMessage("MongoDB initialization failed",
			zap.String("operation", "UpdatePermission"),
			zap.String("permission_id", objectId.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	// 5. Start session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage("MongoDB session start failed",
			zap.String("operation", "UpdatePermission"),
			zap.String("permission_id", objectId.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	// 6. Execute update operation
	err = md.Permission.UpdatePermission(ctx, &session, objectId.String(), updatedPerm)
	if err != nil {
		logger.NewErrMessage("Failed to update permission",
			zap.String("operation", "UpdatePermission"),
			zap.String("permission_id", objectId.String()),
			zap.Error(err),
		)

		return fmt.Errorf("database operation failed: %w", err)
	}

	logger.NewInfoMessage("Successfully updated permission",
		zap.String("operation", "UpdatePermission"),
		zap.String("permission_id", objectId.String()),
	)

	return nil
}
