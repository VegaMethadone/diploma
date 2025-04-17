package permission

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func (r *PermissionMongo) DeletePermission(
	ctx context.Context,
	tx *mongo.Session,
	uuidId string,
) error {
	if tx == nil {
		return errors.New("transaction session is required")
	}

	if uuidId == "" {
		return errors.New("uuid_id cannot be empty")
	}

	exists, err := r.ExistsPermission(ctx, tx, uuidId)
	if err != nil {
		return fmt.Errorf("existence check failed: %w", err)
	}
	if !exists {
		return fmt.Errorf("permission with uuid_id '%s' not found", uuidId)
	}

	err = mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		res, err := r.collection.DeleteOne(sc, bson.M{"uuid_id": uuidId})
		if err != nil {
			return err
		}
		if res.DeletedCount == 0 {
			return errors.New("no documents were deleted")
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}
