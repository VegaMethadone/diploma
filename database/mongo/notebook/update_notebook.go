package notebook

import (
	"context"
	"errors"
	"fmt"
	"labyrinth/notebook/models/journal"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (r *NotebookMongo) UpdateNotebook(
	ctx context.Context,
	tx *mongo.Session,
	uuidId string,
	notebook *journal.Notebook,
) error {
	if tx == nil {
		return errors.New("transaction session is required")
	}
	if uuidId == "" {
		return errors.New("uuidId cannot be empty")
	}
	if notebook == nil {
		return errors.New("notebook cannot be nil")
	}

	update := bson.M{
		"$set": bson.M{
			"version":    notebook.Version,
			"metadata":   notebook.Metadata,
			"blocks":     notebook.Blocks,
			"updated_at": time.Now(), // Автоматическое обновление времени
		},
	}

	// Транзакционная операция
	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		filter := bson.M{"uuid_id": uuidId}
		result := r.collection.FindOneAndUpdate(
			sc,
			filter,
			update,
			options.FindOneAndUpdate().
				SetReturnDocument(options.After), // Возвращаем обновленный документ
		)
		if result.Err() != nil {
			if errors.Is(result.Err(), mongo.ErrNoDocuments) {
				return fmt.Errorf("notebook with uuid_id %s not found", uuidId)
			}
			return fmt.Errorf("failed to update notebook: %w", result.Err())
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to execute notebook update: %w", err)
	}

	return nil
}
