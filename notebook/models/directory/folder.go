package directory

import (
	"time"

	"github.com/google/uuid"
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

func NewDirectory(employeeId, companyId, divisionId, generatedId, parentId uuid.UUID, version string, isPrimary bool, title, description string) Directory {
	return Directory{
		ID:        primitive.NewObjectID(),
		UuidID:    generatedId.String(),
		ParentId:  parentId.String(),
		IsPrimary: isPrimary,
		Version:   version,
		Metadata:  NewMetadata(employeeId, companyId, divisionId, title, description),
		Folders:   []Folder{},
		Files:     []File{},
	}
}

func NewMetadata(employeeId, compaydId, divisionId uuid.UUID, title, description string) Metadata {
	return Metadata{
		CompanyID:   compaydId.String(),
		DivisionID:  divisionId.String(),
		Title:       title,
		Description: description,
		Tags:        []string{},
		Created:     NewTimestamp(time.Now(), employeeId.String()),
		LastUpdate:  NewTimestamp(time.Now(), employeeId.String()),
		Links:       NewDocumentLinks("", "", "", []string{}),
	}
}

func NewDocumentLinks(read, comment, write string, links []string) DocumentLinks {
	return DocumentLinks{
		Read:        read,
		Comment:     comment,
		Write:       write,
		ActiveLinks: links,
	}
}

func NewTimestamp(t time.Time, author string) Timestamp {
	return Timestamp{
		Date:   t,
		Time:   t,
		Author: author,
	}
}

func NewFolder(fileUUID uuid.UUID, title, description string) Folder {
	return Folder{
		FolderID:    primitive.NewObjectID(),
		FolderUUID:  fileUUID.String(),
		Title:       title,
		Description: description,
	}
}

func NewFile(fileUUID uuid.UUID, title, description string) File {
	return File{
		FileID:      primitive.NewObjectID(),
		FileUUID:    fileUUID.String(),
		Title:       title,
		Description: description,
	}
}
