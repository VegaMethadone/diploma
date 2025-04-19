package directory

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Directory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UuidID    string             `bson:"uuid_id"`
	ParentId  string             `bson:"parent_uuid_id"`
	IsPrimary bool               `bson:"isPrimary"`
	Version   string             `bson:"version"`
	Metadata  Metadata           `bson:"metadata"`
	Folders   []Folder           `bson:"folders"`
	Files     []File             `bson:"files"`
}

type Metadata struct {
	CompanyID   string        `bson:"company_id"`
	DivisionID  string        `bson:"division_id"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Tags        []string      `bson:"tags"`
	Created     Timestamp     `bson:"created"`
	LastUpdate  Timestamp     `bson:"last_update"`
	Links       DocumentLinks `bson:"links"`
}

type DocumentLinks struct {
	Read        string   `bson:"read"`
	Comment     string   `bson:"comment"`
	Write       string   `bson:"write"`
	ActiveLinks []string `bson:"active_links"`
}

type Timestamp struct {
	Date   time.Time `bson:"date"`
	Time   time.Time `bson:"time"`
	Author string    `bson:"author"`
}

type Folder struct {
	FolderID    primitive.ObjectID `bson:"folder_id"`
	FolderUUID  string             `bson:"uuid_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
}

type File struct {
	FileID      primitive.ObjectID `bson:"file_id"`
	FileUUID    string             `bson:"uuid_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
}
