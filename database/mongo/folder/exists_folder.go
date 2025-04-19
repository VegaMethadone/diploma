package folder

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (r *FolderMongo) ExistsFolder(
	ctx context.Context,
	sess *mongo.Session,
	folderUUID string,
) (bool, error) {
	if sess == nil {
		return false, errors.New("session is required")
	}
	if folderUUID == "" {
		return false, errors.New("folder UUID cannot be empty")
	}
	var exists bool
	var err error

	err = mongo.WithSession(ctx, *sess, func(sc mongo.SessionContext) error {
		count, err := r.collection.CountDocuments(
			sc,
			bson.M{"uuid_id": folderUUID},
			options.Count().SetLimit(1),
		)
		exists = count > 0
		return err
	})
	if err != nil {
		return false, fmt.Errorf("failed to check folder existence: %w", err)
	}

	return exists, nil
}
