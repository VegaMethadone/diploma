package notebook

import "go.mongodb.org/mongo-driver/mongo"

type NotebookMongo struct {
	collection *mongo.Collection
}

func NewNotebookMongo(db *mongo.Database, collection string) *NotebookMongo {
	return &NotebookMongo{
		collection: db.Collection(collection),
	}
}
