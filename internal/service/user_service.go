package service

import (
	"PullRequestManage/internal/domain"
	"PullRequestManage/internal/repository"
	"context"
)

type UserService struct {
	userRepo repository.UserRepository
	prRepo   repository.PullRequestRepository
}

func NewUserService(userRepo repository.UserRepository, prRepo repository.PullRequestRepository) *UserService {
	return &UserService{userRepo: userRepo, prRepo: prRepo}
}

func (s *UserService) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	user, err := s.userRepo.SetActive(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetReviewsForUser(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.prRepo.FindByReviewer(ctx, userID)
}
