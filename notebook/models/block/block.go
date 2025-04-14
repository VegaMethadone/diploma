package block

import (
	"time"

	"github.com/google/uuid"
)

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
