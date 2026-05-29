package service

import (
	"context"
	"errors"
	"io"
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

func (s *Service) UploadAttachment(ctx context.Context, file io.Reader, fileSize int, fileType string, post *model.Post) (string, error) {
	fileName := randomString(16)
	fileNameUploaded, err := s.fileStorage.Upload(ctx, fileName, file, int64(fileSize), fileType)
	attachment, err := s.attachments.UpdatePostAttachment(ctx, fileNameUploaded, post.ID, fileType)
	if err != nil {
		return "error on changing value picture_path", errors.New("wrong token")
	}
	return attachment.FilePath, err
}
