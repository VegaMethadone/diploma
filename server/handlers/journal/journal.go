package journal

import (
	notebookLogic "labyrinth/notebook/logic"
)

const (
	userIDKey string = "id"
)

var fsl *notebookLogic.FileSystem = notebookLogic.NewFileSystem()

type JournalHandler struct{}

func NewJournalHandler() JournalHandler { return JournalHandler{} }

type notebookRequest struct {
	Name        string `json: "name"`
	Description string `json: "description"`
}
