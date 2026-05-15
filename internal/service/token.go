package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"messenger/internal/model"
	"time"
)

func GenerateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s *Service) MakeToken(ctx context.Context, userId string) (model.Token, error) {
	var token string
	token = GenerateToken()
	return s.tokens.MakeToken(ctx, userId, token)
}

func (s *Service) CheckToken(ctx context.Context, tokenString string, tokenLifeTimeMinutes int) (bool, error) {
	token, err := s.tokens.GetToken(ctx, tokenString)
	if err != nil {
		log.Println("не найден токен " + token.Token)
		return false, nil
	}
	timeDiffMin := int(time.Since(token.CreatedAt).Minutes())
	if timeDiffMin > tokenLifeTimeMinutes {
		log.Println("истек токен " + token.Token)
		return false, nil
	}
	return true, nil
}
