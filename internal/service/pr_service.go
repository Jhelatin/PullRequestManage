package service

import (
	"PullRequestManage/internal/domain"
	"PullRequestManage/internal/repository"
	"context"
	"math/rand"
	"slices"
	"time"
)

type PRService struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
	prRepo   repository.PullRequestRepository
	rand     *rand.Rand
}

func NewPRService(userRepo repository.UserRepository, teamRepo repository.TeamRepository, prRepo repository.PullRequestRepository) *PRService {
	source := rand.NewSource(time.Now().UnixNano())
	return &PRService{
		userRepo: userRepo,
		teamRepo: teamRepo,
		prRepo:   prRepo,
		rand:     rand.New(source),
	}
}

func (s *PRService) CreatePullRequest(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error) {
	existingPR, err := s.prRepo.FindByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if existingPR != nil {
		return nil, ErrPRExists
	}

	author, err := s.userRepo.FindByID(ctx, authorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, ErrAuthorNotFound
	}

	candidates, err := s.teamRepo.GetActiveReviewerCandidates(ctx, author.TeamName, author.ID)
	if err != nil {
		return nil, err
	}

	s.rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	reviewerCount := min(2, len(candidates))
	reviewerIDs := make([]string, 0, reviewerCount)
	for i := 0; i < reviewerCount; i++ {
		reviewerIDs = append(reviewerIDs, candidates[i].ID)
	}

	newPR := domain.PullRequest{
		ID:          prID,
		Name:        prName,
		AuthorID:    authorID,
		Status:      domain.StatusOpen,
		ReviewerIDs: reviewerIDs,
		CreatedAt:   time.Now(),
	}

	if err := s.prRepo.Create(ctx, newPR); err != nil {
		return nil, err
	}
	return &newPR, nil
}

func (s *PRService) MergePullRequest(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.FindByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == domain.StatusMerged {
		return pr, nil
	}

	now := time.Now()
	pr.Status = domain.StatusMerged
	pr.MergedAt = &now

	if err := s.prRepo.Update(ctx, *pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error) {
	pr, err := s.prRepo.FindByID(ctx, prID)
	if err != nil {
		return nil, "", err
	}
	if pr == nil {
		return nil, "", ErrPRNotFound
	}

	if pr.Status == domain.StatusMerged {
		return nil, "", ErrPRMerged
	}

	if !slices.Contains(pr.ReviewerIDs, oldReviewerID) {
		return nil, "", ErrReviewerNotAssigned
	}

	oldReviewer, err := s.userRepo.FindByID(ctx, oldReviewerID)
	if err != nil {
		return nil, "", err
	}
	if oldReviewer == nil {
		return nil, "", ErrUserNotFound
	}

	allMembers, err := s.teamRepo.GetTeamMembers(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, "", err
	}

	excludeSet := make(map[string]struct{})
	excludeSet[pr.AuthorID] = struct{}{}
	for _, id := range pr.ReviewerIDs {
		excludeSet[id] = struct{}{}
	}

	var candidates []domain.User
	for _, member := range allMembers {
		if member.IsActive {
			if _, ok := excludeSet[member.ID]; !ok {
				candidates = append(candidates, member)
			}
		}
	}

	if len(candidates) == 0 {
		return nil, "", ErrNoCandidateForReassign
	}

	newReviewer := candidates[s.rand.Intn(len(candidates))]
	newReviewerIDs := slices.DeleteFunc(pr.ReviewerIDs, func(id string) bool { return id == oldReviewerID })
	pr.ReviewerIDs = append(newReviewerIDs, newReviewer.ID)

	if err := s.prRepo.Update(ctx, *pr); err != nil {
		return nil, "", err
	}
	return pr, newReviewer.ID, nil
}
