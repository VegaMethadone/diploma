package permission

import (
	notebookLogic "labyrinth/notebook/logic"
)

var fsl *notebookLogic.FileSystem = notebookLogic.NewFileSystem()

type PermissionHandlers struct{}

func NewPermissionHandlers() PermissionHandlers { return PermissionHandlers{} }
