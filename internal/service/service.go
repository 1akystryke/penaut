package service

import "messenger/internal/repository"

type Service struct {
	users       repository.UserRepository
	channels    repository.ChannelRepository
	memberships repository.MembershipRepository
	posts       repository.PostRepository
	tokens      repository.TokenRepository
}

func New(
	users repository.UserRepository,
	channels repository.ChannelRepository,
	memberships repository.MembershipRepository,
	posts repository.PostRepository,
	tokens repository.TokenRepository,
) *Service {
	return &Service{
		users:       users,
		channels:    channels,
		memberships: memberships,
		posts:       posts,
		tokens:      tokens,
	}
}
