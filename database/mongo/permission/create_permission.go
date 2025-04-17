package permission

import (
	"context"
	"errors"
	"fmt"
	"labyrinth/notebook/models/permission"

	"go.mongodb.org/mongo-driver/mongo"
)

func (r *PermissionMongo) CreatePermission(
	ctx context.Context,
	tx *mongo.Session,
	permission *permission.Permission,
) error {
	if tx == nil {
		return errors.New("session is required")
	}

	if permission == nil {
		return errors.New("permission cannot be nil")
	}

	if permission.UuidId == "" {
		return errors.New("uuid_id is required")
	}

	if permission.ResourceType == "" {
		return errors.New("resource_type is required")
	}

	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		_, err := r.collection.InsertOne(sc, permission)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return nil
}
