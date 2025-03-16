package block

const (
	TEXT = iota + 1
	TITLE
	TABLE
	IMG
	MATH
	CHEM
)

type Block struct {
	Type     int      `json:"type"`
	ID       int64    `json:"block_id"`
	Metadata Metadata `json:"metadata"`
	Body     any      `json:"body"`
	Comment  Comment  `json:"comment"`
}

type Metadata struct {
	Created TimestampWithAuthor `json:"created"`
	Updated TimestampWithAuthor `json:"updated"`
}

type TimestampWithAuthor struct {
	Date   string `json:"date"`
	Time   string `json:"time"`
	Author string `json:"author"`
}

type Comment struct {
	Metadata TimestampWithAuthor `json:"metadata"`
	Text     string              `json:"text"`
}

func BlockToString(input int) string {
	switch input {
	case 1:
		return "text"
	case 2:
		return "title"
	case 3:
		return "table"
	case 4:
		return "img"
	case 5:
		return "math"
	case 6:
		return "chem"
	default:
		return "[ ERROR ]: block number out of range"
	}
}
