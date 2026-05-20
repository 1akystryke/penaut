package service

import (
	"messenger/internal/repository"
	"messenger/internal/storage/minio"
)

type Service struct {
	users       repository.UserRepository
	channels    repository.ChannelRepository
	memberships repository.MembershipRepository
	posts       repository.PostRepository
	tokens      repository.TokenRepository
	attachments repository.AttachmentRepository
	fileStorage minio.Storage
}

func New(
	users repository.UserRepository,
	channels repository.ChannelRepository,
	memberships repository.MembershipRepository,
	posts repository.PostRepository,
	tokens repository.TokenRepository,
	attachments repository.AttachmentRepository,
	fileStorage minio.Storage,
) *Service {
	return &Service{
		users:       users,
		channels:    channels,
		memberships: memberships,
		posts:       posts,
		tokens:      tokens,
		attachments: attachments,
		fileStorage: fileStorage,
	}
}
