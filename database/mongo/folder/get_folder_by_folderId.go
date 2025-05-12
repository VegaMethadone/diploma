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

func (r *FolderMongo) GetFolderByFolderId(
	ctx context.Context,
	tx *mongo.Session,
	folderId string,
) (*directory.Directory, error) {
	if tx == nil {
		return nil, errors.New("transaction session is required")
	}

	if folderId == "" {
		return nil, errors.New("folderId cannot be empty")
	}

	var results []*directory.Directory
	filter := bson.M{"uuid_id": folderId}
	findOpts := options.Find().SetLimit(1)

	executeQuery := func(ctx context.Context) error {
		cursor, err := r.collection.Find(ctx, filter, findOpts)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &results); err != nil {
			return fmt.Errorf("failed to decode results: %w", err)
		}

		if len(results) == 0 {
			return fmt.Errorf("folder not found")
		}

		return nil
	}

	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		return executeQuery(sc)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	return results[0], nil
}
