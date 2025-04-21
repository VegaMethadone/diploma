package notebook

import (
	"context"
	"errors"
	"fmt"
	"labyrinth/notebook/models/journal"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (r *NotebookMongo) CreateNotebook(
	ctx context.Context,
	tx *mongo.Session,
	notebook *journal.Notebook,
) error {
	if tx == nil {
		return errors.New("transaction session is required")
	}
	if notebook == nil {
		return errors.New("notebook cannot be nil")
	}
	if notebook.UuidID == "" {
		return errors.New("notebook UUID is required")
	}
	if notebook.Metadata.Created.Author == "" {
		return errors.New("author ID is required")
	}
	if notebook.Metadata.Title == "" {
		return errors.New("title is required")
	}

	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		count, err := r.collection.CountDocuments(
			sc,
			bson.M{"uuid": notebook.UuidID},
			options.Count().SetLimit(1),
		)
		if err != nil {
			return fmt.Errorf("failed to check notebook uniqueness: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("notebook with UUID %s already exists", notebook.UuidID)
		}

		_, err = r.collection.InsertOne(sc, notebook)
		if err != nil {
			return fmt.Errorf("failed to insert notebook: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create notebook: %w", err)
	}

	return nil
}
