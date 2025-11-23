package postgres

import (
	"PullRequestManage/internal/domain"
	"PullRequestManage/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgTeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) repository.TeamRepository {
	return &pgTeamRepository{pool: pool}
}

func (r *pgTeamRepository) GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error) {
	query := `
        SELECT user_id, username, team_name, is_active
        FROM users
        WHERE team_name = $1
    `
	rows, err := r.pool.Query(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to query team members: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect user rows: %w", err)
	}

	return users, nil
}

func (r *pgTeamRepository) GetActiveReviewerCandidates(ctx context.Context, teamName string, excludeUserID string) ([]domain.User, error) {
	query := `
        SELECT user_id, username, team_name, is_active
        FROM users
        WHERE team_name = $1 AND is_active = TRUE AND user_id != $2
    `
	rows, err := r.pool.Query(ctx, query, teamName, excludeUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to query active reviewer candidates: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect user rows for candidates: %w", err)
	}

	return users, nil
}

func (r *pgTeamRepository) Create(ctx context.Context, team domain.Team) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO teams (team_name) VALUES ($1)`, team.Name)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	userRepo := NewUserRepository(r.pool)
	if err := userRepo.CreateOrUpdate(ctx, team.Members, team.Name); err != nil {
		return fmt.Errorf("failed to create or update users for team %s: %w", team.Name, err)
	}

	return tx.Commit(ctx)
}

func (r *pgTeamRepository) FindByName(ctx context.Context, name string) (*domain.Team, error) {
	var teamName string
	err := r.pool.QueryRow(ctx, `SELECT team_name FROM teams WHERE team_name = $1`, name).Scan(&teamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find team by name: %w", err)
	}

	members, err := r.GetTeamMembers(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get members for team %s: %w", teamName, err)
	}

	return &domain.Team{
		Name:    teamName,
		Members: members,
	}, nil
}
