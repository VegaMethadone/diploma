package file

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func (f *FileMINIO) DownloadFile(
	ctx context.Context,
	bucketName string,
	objectName string,
	opts minio.GetObjectOptions,
) (*minio.Object, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name cannot be empty")
	}
	if objectName == "" {
		return nil, fmt.Errorf("object name cannot be empty")
	}

	obj, err := f.Client.GetObject(ctx, bucketName, objectName, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	_, err = obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to access object data: %w", err)
	}

	return obj, nil
}
