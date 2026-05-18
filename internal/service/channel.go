package service

import (
	"context"
	"errors"
	"strings"

	"messenger/internal/model"
)

func (s *Service) CreateChannel(ctx context.Context, channelType string, channelName string) (model.Channel, error) {
	channelType = strings.TrimSpace(channelType)
	if channelType != "direct" && channelType != "public" && channelType != "private" {
		return model.Channel{}, errors.New("type must be direct, public, or private")
	}

	return s.channels.CreateChannel(ctx, channelType, channelName)
}

func (s *Service) ListChannels(ctx context.Context) ([]model.Channel, error) {
	return s.channels.ListChannels(ctx)
}
