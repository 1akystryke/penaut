package service

import (
	"context"
	"messenger/internal/model"

	minioSDK "github.com/minio/minio-go/v7"
)

func (s *Service) GetAttachmentsbyPostID(ctx context.Context, postID string) ([]model.Attachment, error) {
	return s.attachments.GetAttachmentsbyPostID(ctx, postID)
}

func (s *Service) GetAttachment(ctx context.Context, attachmentPath string) (*minioSDK.Object, error) {
	file, err := s.fileStorage.GetFile(ctx, attachmentPath)
	return file, err
}
