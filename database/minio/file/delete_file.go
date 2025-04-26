package file

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func (f *FileMINIO) DeleteFile(
	ctx context.Context,
	bucketName string,
	objectName string,
	opts minio.RemoveObjectOptions,
) error {
	if bucketName == "" || objectName == "" {
		return fmt.Errorf("bucket and object names cannot be empty")
	}

	if err := f.client.RemoveObject(ctx, bucketName, objectName, opts); err != nil {
		return fmt.Errorf("failed to remove object: %w", err)
	}

	return nil
}
