package service

import (
	"context"

	"diet-bot/internal/domain/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, userID int64) (int, error) {
	exists, err := s.repo.Exists(ctx, userID)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, nil
	}

	return s.repo.Create(ctx, userID)
}

func (s *UserService) GetAllChats(ctx context.Context) ([]int64, error) {
	return s.repo.GetAllChats(ctx)
}
