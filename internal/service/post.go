package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"messenger/internal/model"
)

func (s *Service) CreatePost(ctx context.Context, channelID string, text string, userToken string) (model.Post, error) {
	channelID = strings.TrimSpace(channelID)
	text = strings.TrimSpace(text)
	if channelID == "" {
		return model.Post{}, errors.New("channel_id is required")
	}
	if text == "" {
		return model.Post{}, errors.New("text is required")
	}
	if userToken == "" {
		return model.Post{}, errors.New("token is required")
	}
	user, err := s.users.GetUserbyToken(ctx, userToken)
	if err != nil {
		log.Println(err)
		log.Println("не найден пользователь по токену " + userToken)
		return model.Post{}, nil
	}
	return s.posts.CreatePost(ctx, channelID, text, user.ID)
}

func (s *Service) ListPosts(ctx context.Context, channelID string) ([]model.Post, error) {
	channelID = strings.TrimSpace(channelID)
	if channelID == "" {
		return nil, errors.New("channel_id is required")
	}

	return s.posts.ListPosts(ctx, channelID)
}
