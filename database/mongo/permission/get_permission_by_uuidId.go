package permission

import (
	"context"
	"errors"
	"fmt"
	"labyrinth/notebook/models/permission"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func (r PermissionMongo) GetPermissionByUuidId(
	ctx context.Context,
	tx *mongo.Session,
	uuidId string,
) (*permission.Permission, error) {
	if tx == nil {
		return nil, errors.New("transaction session is required")
	}

	if uuidId == "" {
		return nil, errors.New("uuid_id cannot be empty")
	}

	filter := bson.M{"uuid_id": uuidId}
	var result permission.Permission

	err := mongo.WithSession(ctx, *tx, func(sc mongo.SessionContext) error {
		return r.collection.FindOne(sc, filter).Decode(&result)
	})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &result, nil
}
