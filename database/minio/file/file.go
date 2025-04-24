package file

import "github.com/minio/minio-go/v7"

type FileMINIO struct {
	Client *minio.Client
}

func NewFileMINIO(client *minio.Client) *FileMINIO {
	return &FileMINIO{
		Client: client,
	}
}
