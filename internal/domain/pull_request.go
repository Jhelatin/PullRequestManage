package domain

import "time"

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID          string     `json:"pull_request_id"`
	Name        string     `json:"pull_request_name"`
	AuthorID    string     `json:"author_id"`
	Status      PRStatus   `json:"status"`
	ReviewerIDs []string   `json:"assigned_reviewers"`
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
	MergedAt    *time.Time `json:"mergedAt,omitempty"`
}
