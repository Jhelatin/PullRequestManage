package service

import (
	"PullRequestManage/internal/domain"
	"PullRequestManage/internal/repository"
	"context"
)

type TeamService struct {
	teamRepo repository.TeamRepository
}

func NewTeamService(teamRepo repository.TeamRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) CreateTeam(ctx context.Context, team domain.Team) (*domain.Team, error) {
	existingTeam, err := s.teamRepo.FindByName(ctx, team.Name)
	if err != nil {
		return nil, err
	}
	if existingTeam != nil {
		return nil, ErrTeamExists
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	team, err := s.teamRepo.FindByName(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}
