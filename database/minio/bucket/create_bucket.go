package bucket

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func (b *BucketMINIO) CreateBucket(
	ctx context.Context,
	bucketName string,
	opts minio.MakeBucketOptions,
) error {
	exists, err := b.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if exists {
		return fmt.Errorf("bucket '%s' already exists", bucketName)
	}

	err = b.Client.MakeBucket(ctx, bucketName, opts)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return nil
}
