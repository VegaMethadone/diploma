package mongo

import (
	"context"
	"fmt"
	"labyrinth/config"
	"labyrinth/notebook/models/permission"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type notebookMongo interface {
	// CreateNotebook

	// UpdateNotebook

	// GetNodebookById

	// GetNotebooksById

	// DeleteNotebook
}

type folderMongo interface {
	// CreateFolder

	// UpdateFolder

	// GetFolderByFolderId

	// GetFoldersByParentId

	// DeleteFolder

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
		uuidId string,
	) (bool, error)
}
type MongoDB struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Folder     *folderMongo
	Notebook   *notebookMongo
	Permission *permissionMongo
}

func NewMongoDB() MongoDB {
	return MongoDB{}
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
