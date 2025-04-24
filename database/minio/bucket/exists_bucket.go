package bucket

import (
	"context"
	"fmt"
)

func (b *BucketMINIO) ExistsBucket(
	ctx context.Context,
	bucketName string,
) (bool, error) {
	if bucketName == "" {
		return false, fmt.Errorf("bucket name cannot be empty")
	}

	exists, err := b.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	return exists, nil
}
