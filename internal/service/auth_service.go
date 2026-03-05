package service

import (
	"context"

	"forms/internal/model"
	"forms/internal/repository"
	"forms/internal/utils"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
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

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repo.Get(ctx, email)
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return "", err
	}

	return utils.GenerateToken(user.ID)
}
