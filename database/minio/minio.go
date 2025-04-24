package minio

import (
	"context"
	"io"
	"labyrinth/database/minio/bucket"
	"labyrinth/database/minio/file"

	"github.com/minio/minio-go/v7"
)

type minioBucket interface {
	// CreateBucket создает бакет
	CreateBucket(
		ctx context.Context,
		bucketName string,
		opts minio.MakeBucketOptions,
	) error

	// ExistsBucket проверяет существование бакета
	ExistsBucket(
		ctx context.Context,
		bucketName string,
	) (bool, error)

	// DeleteExists удаляет бакет
	DeleteForce(
		ctx context.Context,
		bucketName string,
	) error
}

type minioFile interface {
	// UploadFile загружает файл в MinIO
	UploadFile(
		ctx context.Context,
		bucketName string,
		objectName string,
		file io.Reader,
		objectSize int64,
		opts minio.PutObjectOptions,
	) error

	// DownloadFile скачивает файл из MinIO
	DownloadFile(
		ctx context.Context,
		bucketName string,
		objectName string,
		opts minio.GetObjectOptions,
	) (*minio.Object, error)

	// DeleteFile удаляет файл из MinIO
	DeleteFile(
		ctx context.Context,
		bucketName string,
		objectName string,
		opts minio.RemoveObjectOptions,
	) error

	// FileExists проверяет существование файла
	FileExists(
		ctx context.Context,
		bucketName string,
		objectName string,
	) (bool, error)
}

type MinioDB struct {
	Bucket minioBucket
	File   minioFile
}

func NewMinioDB(client *minio.Client) *MinioDB {
	return &MinioDB{
		Bucket: bucket.NewBucketMINIO(client),
		File:   file.NewFileMINIO(client),
	}
}
