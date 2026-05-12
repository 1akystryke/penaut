package service

import (
	"context"
	"errors"
	"strings"

	"messenger/internal/model"
)

func (s *Service) CreateUser(ctx context.Context, name string, email string) (model.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" || email == "" {
		return model.User{}, errors.New("name and email are required")
	}

	return s.users.CreateUser(ctx, name, email)
}

func (s *Service) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.users.ListUsers(ctx)
}
