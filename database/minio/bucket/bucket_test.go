package bucket_test

import (
	"context"
	"fmt"
	"labyrinth/database/minio/bucket"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient *minio.Client
	testBucket  string = "test-bucket"
)

func setupMinio() error {
	var err error
	minioClient, err = minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	return nil
}

func teardownMinio() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	removeOpts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}
	for obj := range minioClient.RemoveObjects(ctx, testBucket,
		minioClient.ListObjects(ctx, testBucket, minio.ListObjectsOptions{}),
		removeOpts,
	) {
		if obj.Err != nil {
			log.Printf("Failed to remove %s: %v", obj.ObjectName, obj.Err)
		}
	}

	minioClient.RemoveBucket(ctx, testBucket)
}

func TestMain(m *testing.M) {
	if err := setupMinio(); err != nil {
		fmt.Printf("MinIO test setup failed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	teardownMinio()
	os.Exit(code)
}

func TestBucketOperations(t *testing.T) {
	ctx := context.Background()
	bucketRepo := bucket.NewBucketMINIO(minioClient)

	t.Run("CreateBucket", func(t *testing.T) {
		err := bucketRepo.CreateBucket(ctx, testBucket, minio.MakeBucketOptions{})
		if err != nil {
			t.Fatalf("CreateBucket failed: %v", err)
		}
	})

	t.Run("ExistsBucket", func(t *testing.T) {
		exists, err := bucketRepo.ExistsBucket(ctx, testBucket)
		if err != nil {
			t.Fatalf("ExistsBucket failed: %v", err)
		}
		if !exists {
			t.Error("Bucket should exist")
		}
	})

	t.Run("DeleteForce", func(t *testing.T) {
		_, _ = minioClient.PutObject(ctx, testBucket, "testfile", strings.NewReader("test"), 4, minio.PutObjectOptions{})

		err := bucketRepo.DeleteForce(ctx, testBucket)
		if err != nil {
			t.Fatalf("DeleteForce failed: %v", err)
		}

		exists, _ := bucketRepo.ExistsBucket(ctx, testBucket)
		if exists {
			t.Error("Bucket should be deleted")
		}
	})
}
