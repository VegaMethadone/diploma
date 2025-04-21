package notebook

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *NotebookMongo) DeleteNotebook(
	ctx context.Context,
	tx *mongo.Session,
	notebookId string,
) error {
	if tx == nil {
		return errors.New("transaction session is required")
	}
	if notebookId == "" {
		return errors.New("notebookId cannot be empty")
	}

	exists, err := r.ExistsNotebook(ctx, tx, notebookId)
	if err != nil {
		return fmt.Errorf("existence check failed: %w", err)
	}
	if !exists {
		return fmt.Errorf("notebook with uuid_id '%s' not found", notebookId)
	}

	err = mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		res, err := r.collection.DeleteOne(sc, bson.M{"uuid_id": notebookId})
		if err != nil {
			return err
		}
		if res.DeletedCount == 0 {
			return errors.New("no notebook were deleted")
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to delete notebook: %w", err)
	}

	return nil
}
