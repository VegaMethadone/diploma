package bucket

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
)

func (b *BucketMINIO) DeleteForce(
	ctx context.Context,
	bucketName string,
) error {
	err := b.Client.RemoveBucket(ctx, bucketName)
	if err == nil {
		return nil
	}

	if !strings.Contains(err.Error(), "bucket not empty") {
		return fmt.Errorf("initial bucket removal failed: %w", err)
	}

	objectsCh := b.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for rErr := range b.Client.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		if rErr.Err != nil {
			return fmt.Errorf("failed to remove object: %w", rErr.Err)
		}
	}

	return b.Client.RemoveBucket(ctx, bucketName)
}
