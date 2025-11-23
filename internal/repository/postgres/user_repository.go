package postgres

import (
	"PullRequestManage/internal/domain"
	"PullRequestManage/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgUserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &pgUserRepository{pool: pool}
}

func (r *pgUserRepository) CreateOrUpdate(ctx context.Context, users []domain.User, teamName string) error {
	batch := &pgx.Batch{}
	query := `
        INSERT INTO users (user_id, username, team_name, is_active)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) DO UPDATE
        SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active;
    `
	for _, user := range users {
		batch.Queue(query, user.ID, user.Username, teamName, user.IsActive)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for range users {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch insert/update for users: %w", err)
		}
	}

	return nil
}

func (r *pgUserRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1`

	var user domain.User
	row := r.pool.QueryRow(ctx, query, userID)

	err := row.Scan(&user.ID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, nil
}

func (r *pgUserRepository) SetActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	query := `
        UPDATE users SET is_active = $1 WHERE user_id = $2
        RETURNING user_id, username, team_name, is_active
    `
	var updatedUser domain.User
	row := r.pool.QueryRow(ctx, query, isActive, userID)

	err := row.Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.TeamName, &updatedUser.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to set active flag for user: %w", err)
	}

	return &updatedUser, nil
}
