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
	Links       Links          `bson:"links"`
}

type DateTimeAuthor struct {
	Date   time.Time `bson:"date"`
	Time   time.Time `bson:"time"`
	Author string    `bson:"author"`
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
	EmployeeId string    `bson:"employee_id"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
	Comment    string    `bson:"comment"`
	SubComment []Comment `bson:"sub_comments"`
}

func NewNotebook(employeeId, companyId, divisionId, generatedId, title, description string) Notebook {
	return Notebook{
		ID:       primitive.NewObjectID(),
		UuidID:   generatedId,
		Version:  "1.0.0",
		Metadata: NewMethadata(employeeId, companyId, divisionId, title, description),
		Blocks:   []Block{},
	}
}

func NewMethadata(employeeId, companyId, divisionId, title, description string) Metadata {
	return Metadata{
		CompanyID:   companyId,
		DivisionID:  divisionId,
		Title:       title,
		Description: description,
		Tags:        []string{},
		Created:     NewDateTimeAuthor(employeeId),
		LastUpdate:  NewDateTimeAuthor(employeeId),
		Links:       Links{},
	}
}

func NewDateTimeAuthor(employeeId string) DateTimeAuthor {
	return DateTimeAuthor{
		Date:   time.Now(),
		Time:   time.Now(),
		Author: employeeId,
	}
}
