package mongo

import (
	"context"
	"fmt"
	"labyrinth/config"
	"labyrinth/database/mongo/folder"
	"labyrinth/database/mongo/notebook"
	mongoPerm "labyrinth/database/mongo/permission"
	"labyrinth/notebook/models/directory"
	"labyrinth/notebook/models/journal"
	"labyrinth/notebook/models/permission"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type notebookMongo interface {
	// CreateNotebook
	CreateNotebook(
		ctx context.Context,
		tx *mongo.Session,
		notebook *journal.Notebook,
	) error

	// UpdateNotebook
	UpdateNotebook(
		ctx context.Context,
		tx *mongo.Session,
		notebookId string,
		notebook *journal.Notebook,
	) error

	// ExistsNotebook
	ExistsNotebook(
		ctx context.Context,
		tx *mongo.Session,
		notebookId string,
	) (bool, error)

	// GetNodebookById
	GetNotebookById(
		ctx context.Context,
		tx *mongo.Session,
		notebookId string,
	) (*journal.Notebook, error)

	// DeleteNotebook
	DeleteNotebook(
		ctx context.Context,
		tx *mongo.Session,
		notebookId string,
	) error
}

type folderMongo interface {
	// CreateFolder создает новую папку
	CreateFolder(
		ctx context.Context,
		tx *mongo.Session,
		folder *directory.Directory,
	) error

	// UpdateFolder обновляет существующую папку
	UpdateFolder(
		ctx context.Context,
		tx *mongo.Session,
		folderId string,
		updateData *directory.Directory,
	) error

	// GetFolderByFolderId возвращает папку по её ID
	// GetFolderByFolderId(
	// 	ctx context.Context,
	// 	tx *mongo.Session,
	// 	folderId string,
	// ) (*directory.Directory, error)

	// GetFoldersByParentId возвращает все папки по ID родительской папки
	GetFoldersByParentId(
		ctx context.Context,
		tx *mongo.Session,
		parentId string,
		opts ...*options.FindOptions,
	) ([]*directory.Directory, error)

	// DeleteFolder удаляет папку по ID
	DeleteFolder(
		ctx context.Context,
		tx *mongo.Session,
		folderId string,
	) error

	// ExistsFolder проверяет существование папки
	ExistsFolder(
		ctx context.Context,
		sess *mongo.Session,
		folderId string,
	) (bool, error)
}

type permissionMongo interface {
	// CreatePermission создает новое разрешение
	CreatePermission(
		ctx context.Context,
		tx *mongo.Session,
		permission *permission.Permission,
	) error

	// UpdatePermission обновляет существующее разрешение
	UpdatePermission(
		ctx context.Context,
		tx *mongo.Session,
		uuidId string,
		updateData *permission.Permission,
	) error

	// GetPermissionByUuidId возвращает разрешение по UUID
	GetPermissionByUuidId(
		ctx context.Context,
		tx *mongo.Session,
		uuidId string,
	) (*permission.Permission, error)

	// DeletePermission удаляет разрешение
	DeletePermission(
		ctx context.Context,
		tx *mongo.Session,
		uuidId string,
	) error

	//  ExistsPermission проверяет сущестует ли  объект в коллекции
	ExistsPermission(
		ctx context.Context,
		tx *mongo.Session,
		uuidId string,
	) (bool, error)
}
type MongoDB struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Folder     folderMongo
	Notebook   notebookMongo
	Permission permissionMongo
}

func NewMongoDB() (*MongoDB, error) {
	client, err := NewConnection()
	if err != nil {
		return nil, err
	}
	db := client.Database("labyrinth")
	return &MongoDB{
		Client:     client,
		Database:   db,
		Folder:     folder.NewFolderMongo(db, "folder"),
		Notebook:   notebook.NewNotebookMongo(db, "notebook"),
		Permission: mongoPerm.NewPermissionMongo(db, "permission"),
	}, nil
}

func NewConnection() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := config.Conf.Mongo.URI + "/?replicaSet=rs0"
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(5 * time.Second).
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}
