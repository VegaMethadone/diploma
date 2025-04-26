package file

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

func (f *FileMINIO) UploadFile(
	ctx context.Context,
	bucketName string,
	objectName string,
	file io.Reader,
	objectSize int64,
	opts minio.PutObjectOptions,
) error {
	if bucketName == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}
	if objectName == "" {
		return fmt.Errorf("object name cannot be empty")
	}
	if file == nil {
		return fmt.Errorf("file reader cannot be nil")
	}

	exists, err := f.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bucket %s does not exist", bucketName)
	}

	_, err = f.client.PutObject(ctx, bucketName, objectName, file, objectSize, opts)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}
