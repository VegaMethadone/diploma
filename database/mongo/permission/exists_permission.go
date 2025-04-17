package permission

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (r *PermissionMongo) ExistsPermission(ctx context.Context, sess *mongo.Session, uuidId string) (bool, error) {
	var exists bool
	var err error
	if sess == nil {
		return false, errors.New("session is required")
	}

	err = mongo.WithSession(ctx, *sess, func(sc mongo.SessionContext) error {
		count, err := r.collection.CountDocuments(
			sc,
			bson.M{"uuid_id": uuidId},
			options.Count().SetLimit(1),
		)
		exists = count > 0
		return err
	})

	return exists, err
}
