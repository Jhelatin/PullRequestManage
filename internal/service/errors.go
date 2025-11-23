package service

import "errors"

var (
	ErrTeamExists             = errors.New("team already exists")
	ErrTeamNotFound           = errors.New("team not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrAuthorNotFound         = errors.New("author not found")
	ErrPRExists               = errors.New("pull request with this id already exists")
	ErrPRNotFound             = errors.New("pull request not found")
	ErrPRMerged               = errors.New("operation is not allowed on a merged pull request")
	ErrReviewerNotAssigned    = errors.New("user is not an assigned reviewer for this pull request")
	ErrNoCandidateForReassign = errors.New("no active replacement candidate available")
)
