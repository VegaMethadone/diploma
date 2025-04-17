package journal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notebook struct {
	ID       primitive.ObjectID `bson:"id"`
	UuidID   string             `bson:"uuid_id"`
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
	Id      int            `bson:"id"`
	Type    string         `bson:"type"`
	Body    map[string]any `bson:"body"`
	Comment []Comment      `bson:"comments"`
}

type Comment struct {
	EmployeeId string     `bson:"employee_id"`
	CreatedAt  time.Time  `bson:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at"`
	Comment    string     `bson:"comment"`
	SubComment []*Comment `bson:"sub_comments"`
}
