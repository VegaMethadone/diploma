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

func (p PermissionMongoLogic) GetPermission(objectId uuid.UUID) (*permission.Permission, error) {
	// 1. Validate input
	if objectId == uuid.Nil {
		logger.NewErrMessage("Invalid permission ID",
			zap.String("operation", "GetPermission"),
			zap.Error(fmt.Errorf("permission ID cannot be nil")),
		)
		return nil, errors.New("invalid permission ID: cannot be nil")
	}

	// 2. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Initialize MongoDB
	md, err := mongo.NewMongoDB()
	if err != nil {
		logger.NewErrMessage("MongoDB initialization failed",
			zap.String("operation", "GetPermission"),
			zap.String("permission_id", objectId.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	// 4. Start session
	session, err := md.Client.StartSession()
	if err != nil {
		logger.NewErrMessage("MongoDB session start failed",
			zap.String("operation", "GetPermission"),
			zap.String("permission_id", objectId.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	// 5. Get permission from MongoDB
	fetchedPermission, err := md.Permission.GetPermissionByUuidId(ctx, &session, objectId.String())
	if err != nil {
		logger.NewErrMessage("Failed to get permission",
			zap.String("operation", "GetPermission"),
			zap.String("permission_id", objectId.String()),
			zap.Error(err),
		)

		return nil, fmt.Errorf("database operation failed: %w", err)
	}

	// 6. Validate fetched permission
	if fetchedPermission == nil {
		logger.NewWarnMessage("Permission not found",
			zap.String("operation", "GetPermission"),
			zap.String("permission_id", objectId.String()),
		)
		return nil, errors.New("permission not found")
	}

	logger.NewInfoMessage("Successfully retrieved permission",
		zap.String("operation", "GetPermission"),
		zap.String("permission_id", objectId.String()),
	)

	return fetchedPermission, nil
}
