package journal

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notebook struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Version  string             `bson:"version"`
	Metadata Metadata           `bson:"metadata"`
	Blocks   []Block            `bson:"blocks"`
}

type Metadata struct {
	CompanyID   string         `bson:"company_id"`
	DivisionID  string         `bson:"division_id"`
	Title       string         `bson:"title"`
	Description string         `bson:"description"`
	Tags        []string       `bson:"tags"`
	Created     DateTimeAuthor `bson:"created"`
	LastUpdate  DateTimeAuthor `bson:"last_update"`
	AccessRules AccessRules    `bson:"access_rules"`
	Links       Links          `bson:"links"`
}

type DateTimeAuthor struct {
	Date   string `bson:"date"`
	Time   string `bson:"time"`
	Author string `bson:"author"`
}
type AccessRules struct {
	AccessDenied  []string `bson:"access_denied"`
	AccessAllowed []string `bson:"access_allowed"`
	ReadOnly      []string `bson:"read_only"`
	CommentOnly   []string `bson:"comment_only"`
	AccessLevel   string   `bson:"access_level"`
}

type Links struct {
	Read        string   `bson:"read"`
	Comment     string   `bson:"comment"`
	Write       string   `bson:"write"`
	ActiveLinks []string `bson:"active_links"`
}

type Block struct {
	BlockType    string
	BlockId      int
	BlockBody    any
	BlockComment []Comment
}

type Comment struct {
	EmployeeId uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Comment    string
	SubComment []*Comment
}
