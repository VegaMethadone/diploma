package file

import "github.com/minio/minio-go/v7"

type FileMINIO struct {
	client *minio.Client
}

func NewFileMINIO(client *minio.Client) *FileMINIO {
	return &FileMINIO{
		client: client,
	}
}
