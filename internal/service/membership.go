package service

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/model"
	"strings"
)

func (s *Service) AddMember(ctx context.Context, userID string, channelID string, requestorToken string) (model.Membership, error) {
	userID = strings.TrimSpace(userID)
	channelID = strings.TrimSpace(channelID)
	requestor, err := s.users.GetUserbyToken(ctx, requestorToken)
	if err != nil {
		fmt.Println(err)
		return model.Membership{}, errors.New("wrong token")
	}

	currentMembersList, err := s.memberships.ListMembers(ctx, channelID)

	loopBraker := false
	for _, v := range currentMembersList {
		if v.ID == requestor.ID {
			loopBraker = true
		}
	}
	if !loopBraker {
		return model.Membership{}, errors.New("requestor is not in channel")
	}

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
