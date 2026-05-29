package repository

import (
	"context"

	"messenger/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, name string, email string, passwordHash string) (model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	CheckPassword(ctx context.Context, email string, passwordHash string) (*model.User, error)
	GetUserbyToken(ctx context.Context, token string) (*model.User, error)
	GetUserbyID(ctx context.Context, userID string) (*model.User, error)
}

type ChannelRepository interface {
	CreateChannel(ctx context.Context, channelType string, channelName string) (model.Channel, error)
	ListChannels(ctx context.Context, userId string) ([]model.Channel, error)
	GetChannelbyID(ctx context.Context, channelID string) (*model.Channel, error)
	UpdateChannelPic(ctx context.Context, picName string, channel *model.Channel) (bool, error)
}

type MembershipRepository interface {
	AddMember(ctx context.Context, userID string, channelID string) (model.Membership, error)
	ListMembers(ctx context.Context, channelID string) ([]model.User, error)
}

type PostRepository interface {
	CreatePost(ctx context.Context, channelID string, text string, userID string) (model.Post, error)
	ListPosts(ctx context.Context, channelID string) ([]model.Post, error)
}

type TokenRepository interface {
	MakeToken(ctx context.Context, userID string, token string) (model.Token, error)
	GetToken(ctx context.Context, token string) (*model.Token, error)
}

type AttachmentRepository interface {
	GetAttachmentsbyPostID(ctx context.Context, postID string) ([]model.Attachment, error)
	UpdatePostAttachment(ctx context.Context, picName string, postID string, fileType string) (*model.Attachment, error)
}
