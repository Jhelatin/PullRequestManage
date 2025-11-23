package repository

import (
	"PullRequestManage/internal/domain"
	"context"
)

type UserRepository interface {
	CreateOrUpdate(ctx context.Context, users []domain.User, teamName string) error
	FindByID(ctx context.Context, userID string) (*domain.User, error)
	SetActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team domain.Team) error
	FindByName(ctx context.Context, name string) (*domain.Team, error)
	GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error)
	GetActiveReviewerCandidates(ctx context.Context, teamName string, excludeUserID string) ([]domain.User, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr domain.PullRequest) error
	Update(ctx context.Context, pr domain.PullRequest) error
	FindByID(ctx context.Context, prID string) (*domain.PullRequest, error)
	FindByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error)
}
