package folder

import (
	"context"
	"errors"
	"fmt"
	"labyrinth/notebook/models/directory"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (r *FolderMongo) UpdateFolder(
	ctx context.Context,
	tx *mongo.Session,
	folderId string,
	updateData *directory.Directory,
) error {
	if tx == nil {
		return errors.New("transaction session is required")
	}

	if folderId == "" {
		return errors.New("folderId cannot be empty")
	}

	if updateData == nil {
		return errors.New("updateData cannot be nil")
	}

	filter := bson.M{"uuid_id": folderId}
	update := bson.M{
		"$set": bson.M{
			"isPrimary": updateData.IsPrimary,
			"version":   updateData.Version,
			"metadata":  updateData.Metadata,
			"folders":   updateData.Folders,
			"files":     updateData.Files,
		},
	}

	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		count, err := r.collection.CountDocuments(
			sc,
			filter,
			options.Count().SetLimit(1),
		)
		if err != nil {
			return fmt.Errorf("failed to check folder existence: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("folder with id %s not found", folderId)
		}

		result, err := r.collection.UpdateOne(
			sc,
			filter,
			update,
			options.Update().SetUpsert(false),
		)
		if err != nil {
			return fmt.Errorf("failed to update folder: %w", err)
		}

		if result.MatchedCount == 0 {
			return fmt.Errorf("no folder matched for update with id %s", folderId)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to execute folder update transaction: %w", err)
	}

	return nil
}
