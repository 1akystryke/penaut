package repository

import (
	"context"

	"messenger/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, name string, email string, passwordHash string) (model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
}

type ChannelRepository interface {
	CreateChannel(ctx context.Context, channelType string) (model.Channel, error)
	ListChannels(ctx context.Context) ([]model.Channel, error)
}

type MembershipRepository interface {
	AddMember(ctx context.Context, userID string, channelID string) (model.Membership, error)
	ListMembers(ctx context.Context, channelID string) ([]model.User, error)
}

type PostRepository interface {
	CreatePost(ctx context.Context, channelID string, text string, userID string) (model.Post, error)
	ListPosts(ctx context.Context, channelID string) ([]model.Post, error)
}
