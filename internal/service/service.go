package service

import "messenger/internal/repository"

type Service struct {
	users       repository.UserRepository
	channels    repository.ChannelRepository
	memberships repository.MembershipRepository
	posts       repository.PostRepository
}

func New(
	users repository.UserRepository,
	channels repository.ChannelRepository,
	memberships repository.MembershipRepository,
	posts repository.PostRepository,
) *Service {
	return &Service{
		users:       users,
		channels:    channels,
		memberships: memberships,
		posts:       posts,
	}
}
