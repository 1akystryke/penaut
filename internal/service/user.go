package service

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/model"
	"strings"

	minioSDK "github.com/minio/minio-go/v7"
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

func (s *Service) GetUserPic(ctx context.Context, user model.User) (*minioSDK.Object, error) {
	file, err := s.fileStorage.GetFile(ctx, user.PicturePath)
	return file, err

}

func (s *Service) GetUserbyID(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return &model.User{}, errors.New("userID required")
	}
	user, err := s.users.GetUserbyID(ctx, userID)
	if err != nil {
		fmt.Println(err)
		return &model.User{}, errors.New("wrong token")
	}
	return user, err
}

func (s *Service) CheckPassword(ctx context.Context, email string, passwordHash string) (*model.User, error) {
	if email == "" || passwordHash == "" {
		return &model.User{}, errors.New("pwd and email are required")
	}
	user, err := s.users.CheckPassword(ctx, email, passwordHash)
	if err != nil {
		fmt.Println(err)
		return &model.User{}, errors.New("wrong password")
	}
	return user, err
}

func (s *Service) GetUserbyToken(ctx context.Context, token string) (*model.User, error) {
	if token == "" {
		return &model.User{}, errors.New("token required")
	}
	user, err := s.users.GetUserbyToken(ctx, token)
	if err != nil {
		fmt.Println(err)
		return &model.User{}, errors.New("wrong token")
	}
	return user, err
}
