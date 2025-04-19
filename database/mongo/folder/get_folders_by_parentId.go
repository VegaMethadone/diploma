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

func (r *FolderMongo) GetFoldersByParentId(
	ctx context.Context,
	tx *mongo.Session,
	parentId string,
	opts ...*options.FindOptions,
) ([]*directory.Directory, error) {
	if tx == nil {
		return nil, errors.New("transaction session is required")
	}

	if parentId == "" {
		return nil, errors.New("parentId cannot be empty")
	}

	filter := bson.M{
		"parent_uuid_id": parentId,
	}

	var results []*directory.Directory
	var findOpts *options.FindOptions

	if len(opts) > 0 {
		findOpts = opts[0]
	}

	executeQuery := func(ctx context.Context) error {
		cursor, err := r.collection.Find(ctx, filter, findOpts)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &results); err != nil {
			return fmt.Errorf("failed to decode results: %w", err)
		}

		return nil
	}

	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		return executeQuery(sc)
	})
	if err != nil {
		return nil, fmt.Errorf("transactional query failed: %w", err)
	}

	return results, nil
}
