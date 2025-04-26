package file_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"labyrinth/database/minio/bucket"
	"labyrinth/database/minio/file"
	"log"
	"os"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient *minio.Client
	testBucket  string = "test-bucket"
	testFile    string = "aGVsbG93b3JsZA=="
	fileName    string = "helloworld"
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

func TestFileOperations(t *testing.T) {
	ctx, err := context.WithTimeout(context.Background(), 5*time.Second)
	if err != nil {
		t.Errorf("Failed to setup context: %v", err)
	}
	bucketRepo := bucket.NewBucketMINIO(minioClient)
	fileRepo := file.NewFileMINIO(minioClient)

	t.Run("CreateBucket", func(t *testing.T) {
		err := bucketRepo.CreateBucket(ctx, testBucket, minio.MakeBucketOptions{})
		if err != nil {
			t.Fatalf("CreateBucket failed: %v", err)
		}
	})

	t.Run("UploadFile", func(t *testing.T) {
		reader := bytes.NewBufferString(testFile)
		err := fileRepo.UploadFile(ctx, testBucket, fileName, reader, int64(len(testFile)), minio.PutObjectOptions{})
		if err != nil {
			t.Fatalf("UploadFile failed: %v", err)
		}
	})

	t.Run("FileExists", func(t *testing.T) {
		exists, err := fileRepo.FileExists(ctx, testBucket, fileName)
		if err != nil {
			t.Fatalf("FileExists failed: %v", err)
		}

		if !exists {
			t.Errorf("Expected exists = true, got exists = false")
		}
	})

	t.Run("DownloadFile", func(t *testing.T) {
		obj, err := fileRepo.DownloadFile(ctx, testBucket, fileName, minio.GetObjectOptions{})
		if err != nil {
			t.Fatalf("DownloadFile failed: %v", err)
		}
		defer obj.Close()

		data, err := io.ReadAll(obj)
		if err != nil {
			t.Fatalf("Failed to copy data: %v", err)
		}

		if string(data) != testFile {
			t.Errorf("Expected %s, but got %s\n", testFile, string(data))
		}
	})

	t.Run("DeleteFile", func(t *testing.T) {
		err := fileRepo.DeleteFile(ctx, testBucket, fileName, minio.RemoveObjectOptions{})
		if err != nil {
			t.Fatalf("DeleteFile failed: %v", err)
		}
	})
}
