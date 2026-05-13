package service

import (
	"context"
	"errors"
	"strings"

	"messenger/internal/model"
)

func (s *Service) CreateUser(ctx context.Context, name string, email string, passwordHash string) (model.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" || email == "" {
		return model.User{}, errors.New("name and email are required")
	}

	return s.users.CreateUser(ctx, name, email, passwordHash)
}

func (s *Service) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.users.ListUsers(ctx)
}

func (s *Service) CheckPassword(ctx context.Context, email string, passwordHash string) (*model.User, error) {
	if email == "" || passwordHash == "" {
		return &model.User{}, errors.New("pwd and email are required")
	}
	user, err := s.users.CheckPassword(ctx, email, passwordHash)
	if err != nil {
		return &model.User{}, errors.New("wrong password")
	}
	return user, err
}
