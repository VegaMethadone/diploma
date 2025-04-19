package folder

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func (r *FolderMongo) DeleteFolder(
	ctx context.Context,
	tx *mongo.Session,
	folderId string,
) error {
	if tx == nil {
		return errors.New("transaction session is required")
	}
	if folderId == "" {
		return errors.New("folderId cannot be empty")
	}

	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		// Выполняем удаление
		filter := bson.M{"uuid_id": folderId}
		result, err := r.collection.DeleteOne(sc, filter)
		if err != nil {
			return fmt.Errorf("failed to delete folder: %w", err)
		}
		if result.DeletedCount == 0 {
			return fmt.Errorf("no folder was deleted (uuid: %s)", folderId)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to execute delete operation: %w", err)
	}

	return nil
}
