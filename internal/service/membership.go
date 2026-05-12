package service

import (
	"context"
	"errors"
	"strings"

	"messenger/internal/model"
)

func (s *Service) AddMember(ctx context.Context, userID string, channelID string) (model.Membership, error) {
	userID = strings.TrimSpace(userID)
	channelID = strings.TrimSpace(channelID)
	if userID == "" {
		return model.Membership{}, errors.New("user_id is required")
	}
	if channelID == "" {
		return model.Membership{}, errors.New("channel_id is required")
	}

	return s.memberships.AddMember(ctx, userID, channelID)
}

func (s *Service) ListMembers(ctx context.Context, channelID string) ([]model.User, error) {
	channelID = strings.TrimSpace(channelID)
	if channelID == "" {
		return nil, errors.New("channel_id is required")
	}

	return s.memberships.ListMembers(ctx, channelID)
}
