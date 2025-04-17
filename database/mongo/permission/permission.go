package permission

import "go.mongodb.org/mongo-driver/mongo"

type PermissionMongo struct {
	collection *mongo.Collection
}

func NewPermissionMongo(db *mongo.Database, collection string) *PermissionMongo {
	return &PermissionMongo{
		collection: db.Collection(collection),
	}
}
