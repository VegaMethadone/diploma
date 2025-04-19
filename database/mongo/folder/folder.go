package folder

import "go.mongodb.org/mongo-driver/mongo"

type FolderMongo struct {
	collection *mongo.Collection
}

func NewFolderMongo(db *mongo.Database, collection string) *FolderMongo {
	return &FolderMongo{
		collection: db.Collection(collection),
	}
}
