package permission

import (
	"context"
	"errors"
	"fmt"
	"labyrinth/notebook/models/permission"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (r *PermissionMongo) UpdatePermission(
	ctx context.Context,
	tx *mongo.Session,
	uuidId string,
	updateData *permission.Permission,
) error {
	if tx == nil {
		return errors.New("transaction session is required")
	}

	if uuidId == "" {
		return errors.New("uuid_id cannot be empty")
	}

	if updateData == nil {
		return errors.New("update data cannot be nil")
	}

	update := bson.M{
		"$set": updateData,
	}

	var result *mongo.UpdateResult
	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		var err error
		result, err = r.collection.UpdateOne(
			sc,
			bson.M{"uuid_id": uuidId},
			update,
			options.Update().SetUpsert(false),
		)
		return err
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("duplicate key violation: %w", err)
		}
		return fmt.Errorf("failed to update permission: %w", err)
	}

	switch {
	case result.MatchedCount == 0:
		return fmt.Errorf("permission with uuid_id '%s' not found", uuidId)
	case result.ModifiedCount == 0:
		return fmt.Errorf("permission data not modified (no changes)")
	}

	return nil
}
