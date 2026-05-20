package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	minioSDK "github.com/minio/minio-go/v7"
)

func (s *Storage) Upload(
	ctx context.Context,
	objectName string,
	reader io.Reader,
	size int64,
	contentType string,
) (string, error) {

	_, err := s.client.PutObject(
		ctx,
		s.bucket,
		objectName,
		reader,
		size,
		minioSDK.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (s *Storage) Delete(
	ctx context.Context,
	objectName string,
) error {

	return s.client.RemoveObject(
		ctx,
		s.bucket,
		objectName,
		minioSDK.RemoveObjectOptions{},
	)
}

func (s *Storage) GetURL(
	ctx context.Context,
	objectName string,
) (string, error) {

	url, err := s.client.PresignedGetObject(
		ctx,
		s.bucket,
		objectName,
		time.Hour,
		nil,
	)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (s *Storage) List(
	ctx context.Context,
) error {

	for object := range s.client.ListObjects(
		ctx,
		s.bucket,
		minioSDK.ListObjectsOptions{
			Recursive: true,
		},
	) {

		if object.Err != nil {
			return object.Err
		}

		fmt.Println(object.Key, object.Size)
	}

	return nil
}

func (s *Storage) GetFile(ctx context.Context, filename string) (*minioSDK.Object, error) {
	log.Println("получаем файл" + filename)
	obj, err := s.client.GetObject(
		ctx,
		"files",
		filename,
		minioSDK.GetObjectOptions{},
	)

	return obj, err
}
