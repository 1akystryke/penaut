package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"messenger/internal/model"
)

func (s *Service) CreateChannel(ctx context.Context, channelType string, channelName string, authorToken string) (model.Channel, error) {
	channelType = strings.TrimSpace(channelType)
	if channelType != "public" && channelType != "private" {
		return model.Channel{}, errors.New("type must be direct, public, or private")
	}
	user, err := s.users.GetUserbyToken(ctx, authorToken)
	channel, err := s.channels.CreateChannel(ctx, channelType, channelName)
	s.memberships.AddMember(ctx, user.ID, channel.ID)
	return channel, err
}

func (s *Service) CreateDirect(ctx context.Context, UserID string, authorToken string) (model.Channel, error) {

	author, err := s.users.GetUserbyToken(ctx, authorToken)

	usersList, err := s.users.ListUsers(ctx)

	loopBraker := false
	partner := model.User{}
	for _, v := range usersList {
		if v.ID == UserID {
			loopBraker = true
			partner = v
		}
	}
	if !loopBraker {
		return model.Channel{}, errors.New("error, user not found")
	}
	channelName := author.ID + "_" + partner.ID
	channel, err := s.channels.CreateChannel(ctx, "direct", channelName)
	s.memberships.AddMember(ctx, author.ID, channel.ID)
	s.memberships.AddMember(ctx, partner.ID, channel.ID)
	log.Println("создали директ")
	return channel, err
}

func (s *Service) ListChannels(ctx context.Context, requestorToken string) ([]model.Channel, error) {
	user, _ := s.users.GetUserbyToken(ctx, requestorToken)
	return s.channels.ListChannels(ctx, user.ID)
}
