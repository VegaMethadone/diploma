package folder

import (
	"context"
	"errors"
	"fmt"
	"labyrinth/notebook/models/directory"

	"go.mongodb.org/mongo-driver/mongo"
)

func (r *FolderMongo) CreateFolder(
	ctx context.Context,
	tx *mongo.Session,
	directory *directory.Directory,
) error {
	if tx == nil {
		return errors.New("transaction session is required")
	}

	if directory == nil {
		return errors.New("directory cannot be nil")
	}

	if directory.UuidID == "" {
		return errors.New("uuid_id is required")
	}

	if directory.Metadata.CompanyID == "" {
		return errors.New("company_id is required in metadata")
	}

	if directory.Metadata.Title == "" {
		return errors.New("title is required in metadata")
	}

	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {

		_, err := r.collection.InsertOne(sc, directory)
		if err != nil {
			return fmt.Errorf("failed to insert directory: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}
