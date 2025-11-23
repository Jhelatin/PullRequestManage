package postgres

import (
	"PullRequestManage/internal/domain"
	"PullRequestManage/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgPullRequestRepository struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepository(pool *pgxpool.Pool) repository.PullRequestRepository {
	return &pgPullRequestRepository{pool: pool}
}

func (r *pgPullRequestRepository) Create(ctx context.Context, pr domain.PullRequest) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	prQuery := `
        INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = tx.Exec(ctx, prQuery, pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	if len(pr.ReviewerIDs) > 0 {
		reviewerQuery := `INSERT INTO reviewers (pull_request_id, user_id) VALUES ($1, $2)`
		for _, reviewerID := range pr.ReviewerIDs {
			_, err = tx.Exec(ctx, reviewerQuery, pr.ID, reviewerID)
			if err != nil {
				return fmt.Errorf("failed to assign reviewer %s to PR %s: %w", reviewerID, pr.ID, err)
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *pgPullRequestRepository) FindByID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	query := `
        SELECT
            pr.pull_request_id,
            pr.pull_request_name,
            pr.author_id,
            pr.status,
            pr.created_at,
            pr.merged_at,
            COALESCE(array_agg(r.user_id) FILTER (WHERE r.user_id IS NOT NULL), '{}') as assigned_reviewers
        FROM pull_requests pr
        LEFT JOIN reviewers r ON pr.pull_request_id = r.pull_request_id
        WHERE pr.pull_request_id = $1
        GROUP BY pr.pull_request_id
    `
	var pr domain.PullRequest
	row := r.pool.QueryRow(ctx, query, prID)

	err := row.Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
		&pr.ReviewerIDs,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pr by id: %w", err)
	}

	return &pr, nil
}

func (r *pgPullRequestRepository) FindByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	query := `
        SELECT
            pr.pull_request_id,
            pr.pull_request_name,
            pr.author_id,
            pr.status
        FROM pull_requests pr
        JOIN reviewers r ON pr.pull_request_id = r.pull_request_id
        WHERE r.user_id = $1
        ORDER BY pr.created_at DESC
    `
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query prs by reviewer: %w", err)
	}
	defer rows.Close()

	prs, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.PullRequest])
	if err != nil {
		return nil, fmt.Errorf("failed to collect pr rows by reviewer: %w", err)
	}

	return prs, nil
}
