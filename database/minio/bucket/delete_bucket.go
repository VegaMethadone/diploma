package bucket

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func (b *BucketMINIO) DeleteForce(
	ctx context.Context,
	bucketName string,
) error {
	objectsCh := b.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for err := range b.client.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		if err.Err != nil {
			return fmt.Errorf("failed to remove object %s: %w", err.ObjectName, err.Err)
		}
	}

	if err := b.client.RemoveBucket(ctx, bucketName); err != nil {
		return fmt.Errorf("failed to remove bucket: %w", err)
	}

	return nil
}
