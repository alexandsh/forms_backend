package service

import (
	"context"
	"errors"
	"time"

	"forms/internal/model"
	"forms/internal/repository"
	"forms/internal/utils"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	repo  *repository.UserRepository
	redis *redis.Client
}

func NewAuthService(repo *repository.UserRepository, redis *redis.Client) *AuthService {
	return &AuthService{
		repo:  repo,
		redis: redis,
	}
}

func (s *AuthService) Register(ctx context.Context, email string, password string) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := model.User{
		Email:        email,
		PasswordHash: hash,
	}

	return s.repo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, string, error) {
	user, err := s.repo.Get(ctx, email)
	if err != nil {
		return "", "", err
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken := utils.GenerateRefreshToken()

	err = s.redis.Set(ctx, refreshToken, user.ID, 7*24*time.Hour).Err()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userID, err := s.redis.Get(ctx, refreshToken).Int()
	if err != nil {
		return "", err
	}

	return utils.GenerateAccessToken(userID)
}
