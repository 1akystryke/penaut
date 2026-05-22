package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"messenger/internal/model"
	"strings"
	"time"

	minioSDK "github.com/minio/minio-go/v7"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}

	return string(b)
}

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

func (s *Service) GetChannelPic(ctx context.Context, channel model.Channel) (*minioSDK.Object, error) {
	file, err := s.fileStorage.GetFile(ctx, channel.PicturePath)
	return file, err

}

func (s *Service) UploadChannelPic(ctx context.Context, channel *model.Channel, file io.Reader, fileSize int, fileType string) (string, error) {
	fileName := randomString(16)
	//file, err := (, channel.PicturePath)
	fileNameUploaded, err := s.fileStorage.Upload(ctx, fileName, file, int64(fileSize), fileType)
	isOk, err := s.channels.UpdateChannelPic(ctx, fileNameUploaded, channel)
	if !isOk {
		return "error on changing value picture_path", errors.New("wrong token")
	}
	return fileNameUploaded, err

}

func (s *Service) GetChannelbyID(ctx context.Context, channelID string) (*model.Channel, error) {
	if channelID == "" {
		return &model.Channel{}, errors.New("channelID required")
	}
	user, err := s.channels.GetChannelbyID(ctx, channelID)
	if err != nil {
		fmt.Println(err)
		return &model.Channel{}, errors.New("wrong token")
	}
	return user, err
}
