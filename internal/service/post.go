package service

import (
	"context"
	"errors"
	"strings"

	"messenger/internal/model"
)

func (s *Service) CreatePost(ctx context.Context, channelID string, text string, userID string) (model.Post, error) {
	channelID = strings.TrimSpace(channelID)
	text = strings.TrimSpace(text)
	if channelID == "" {
		return model.Post{}, errors.New("channel_id is required")
	}
	if text == "" {
		return model.Post{}, errors.New("text is required")
	}

	return s.posts.CreatePost(ctx, channelID, text, userID)
}

func (s *Service) ListPosts(ctx context.Context, channelID string) ([]model.Post, error) {
	channelID = strings.TrimSpace(channelID)
	if channelID == "" {
		return nil, errors.New("channel_id is required")
	}

	return s.posts.ListPosts(ctx, channelID)
}
