package file

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func (f *FileMINIO) FileExists(
	ctx context.Context,
	bucketName string,
	objectName string,
) (bool, error) {
	if bucketName == "" {
		return false, fmt.Errorf("bucket name cannot be empty")
	}
	if objectName == "" {
		return false, fmt.Errorf("object name cannot be empty")
	}

	_, err := f.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if errResp := minio.ToErrorResponse(err); errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}
