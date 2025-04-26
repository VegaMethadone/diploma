package bucket

import "github.com/minio/minio-go/v7"

type BucketMINIO struct {
	client *minio.Client
}

func NewBucketMINIO(client *minio.Client) *BucketMINIO {
	return &BucketMINIO{
		client: client,
	}
}
