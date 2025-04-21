package notebook

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (r *NotebookMongo) ExistsNotebook(
	ctx context.Context,
	tx *mongo.Session,
	notebookId string,
) (bool, error) {
	if tx == nil {
		return false, errors.New("transaction session is required")
	}
	if notebookId == "" {
		return false, errors.New("notebookId cannot be empty")
	}

	var exists bool
	var err error
	err = mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		count, err := r.collection.CountDocuments(
			sc,
			bson.M{"uuid_id": notebookId},
			options.Count().SetLimit(1),
		)
		exists = count > 0
		return err
	})
	if err != nil {
		return false, fmt.Errorf("failed to check notebook existence: %w", err)
	}

	return exists, nil
}
