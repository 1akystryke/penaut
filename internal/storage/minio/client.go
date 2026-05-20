package minio

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client *minio.Client
	bucket string
}

func New(
	endpoint string,
	accessKey string,
	secretKey string,
	bucket string,
	useSSL bool,
) *Storage {

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{
		client: client,
		bucket: bucket,
	}
}
