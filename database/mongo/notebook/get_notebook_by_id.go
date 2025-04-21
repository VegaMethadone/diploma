package notebook

import (
	"context"
	"errors"
	"fmt"
	"labyrinth/notebook/models/journal"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func (r *NotebookMongo) GetNotebookById(
	ctx context.Context,
	tx *mongo.Session,
	notebookId string,
) (*journal.Notebook, error) {
	if tx == nil {
		return nil, errors.New("transaction session is required")
	}
	if notebookId == "" {
		return nil, errors.New("notebookId cannot be empty")
	}

	filter := bson.M{"uuid_id": notebookId}
	var notebook journal.Notebook

	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		return r.collection.FindOne(sc, filter).Decode(&notebook)
	})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("notebook not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &notebook, nil
}
