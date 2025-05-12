package notebookLogic

import (
	folderLogic "labyrinth/notebook/logic/folder"
	notebookLogic "labyrinth/notebook/logic/notebook"
	permissionLogic "labyrinth/notebook/logic/permission"
	"labyrinth/notebook/models/directory"
	"labyrinth/notebook/models/journal"
	"labyrinth/notebook/models/permission"

	"github.com/google/uuid"
)

type notebookInterface interface {
	NewNotebook(employeeId, companyId, divisionId uuid.UUID, title, description string) error
	GetNotebook(notebookId uuid.UUID) (*journal.Notebook, error)
	UpdateNotebook(notebookId uuid.UUID, updatedNotebook *journal.Notebook) error
	DeleteNotebook(notebookId uuid.UUID) error
}

type directoryInterface interface {
	CreateFolder(employeeId, companyId, divisionId, parentId uuid.UUID, isPrimary bool, title, description string) error
	GetFolder(folderId uuid.UUID) (*directory.Directory, error)
	UpdateFolder(folderId uuid.UUID, dir *directory.Directory) error
	DeleteFolder(folderId uuid.UUID) error
}
type permissionInterface interface {
	GetPermission(objectId uuid.UUID) (*permission.Permission, error)
	UpdatePermission(objectId uuid.UUID, updatedPerm *permission.Permission) error
}
type FileSystem struct {
	Folder     directoryInterface
	File       notebookInterface
	Permission permissionInterface
}

func NewFileSystem() *FileSystem {
	return &FileSystem{
		Folder:     folderLogic.NewFolderMongoLogic(),
		File:       notebookLogic.NewNotebookMongoLogic(),
		Permission: permissionLogic.NewPermissionMongoLogic(),
	}
}
