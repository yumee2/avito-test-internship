package pullrequest

import (
	"avito-test/internal/domain/entities"
	"errors"
	"math/rand"
)

type PullRequestService struct {
	pullRequestRepo pullRequestRepo
	userRepo        userRepository
}

func NewPullRequestService(pullRequestRepo pullRequestRepo, userRepo userRepository) *PullRequestService {
	return &PullRequestService{pullRequestRepo: pullRequestRepo, userRepo: userRepo}
}

func (s *PullRequestService) CreatePR(prId string, prName string, authorId string) (*entities.PullRequest, error) {
	pr_reviewers, err := s.userRepo.GetUserFromTheSameTeam(authorId)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}

	pullRequest := &entities.PullRequest{
		PullRequestID:   prId,
		PullRequestName: prName,
		AuthorID:        authorId,
		Status:          string(entities.StatusOpen),
	}

	if len(pr_reviewers) >= 2 {
		idx1 := rand.Intn(len(pr_reviewers))
		idx2 := rand.Intn(len(pr_reviewers))
		for idx1 == idx2 {
			idx2 = rand.Intn(len(pr_reviewers))
		}

		pullRequest.AssignedReviwers = append(pullRequest.AssignedReviwers, pr_reviewers[idx1])
		pullRequest.AssignedReviwers = append(pullRequest.AssignedReviwers, pr_reviewers[idx2])
	} else if len(pr_reviewers) == 1 {
		pullRequest.AssignedReviwers = append(pullRequest.AssignedReviwers, pr_reviewers[0])
	}

	pullRequest, err = s.pullRequestRepo.CreatePR(pullRequest)
	if err != nil {
		if errors.Is(err, entities.ErrDuplicate) {
			return nil, entities.ErrDuplicate
		} else if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.ErrNotFound
		}

		return nil, err
	}

	return pullRequest, nil

}

func (s *PullRequestService) MergePr(pullRequestID string) (*entities.PullRequest, error) {
	pullRequest, err := s.pullRequestRepo.MergePr(pullRequestID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.ErrNotFound
		}

		return nil, err
	}

	return pullRequest, nil
}

func (s *PullRequestService) Reassign(pullRequestID string, oldUserID string) (*entities.PullRequest, error) {
	pullRequest, err := s.pullRequestRepo.GetPR(pullRequestID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}

	if pullRequest.Status == string(entities.StatusMerged) {
		return nil, entities.ErrPrMergerd
	}

	for i, val := range pullRequest.AssignedReviwers {
		if val.UserID == oldUserID {
			pr_reviewers, err := s.userRepo.GetUserForReassign(pullRequest.AssignedReviwers, pullRequest.AuthorID)
			if err != nil {
				if errors.Is(err, entities.ErrNotFound) {
					return nil, entities.ErrNotFound
				}
				return nil, err
			}
			if len(pr_reviewers) == 0 {
				return nil, entities.ErrNoCandidates
			}
			idx := rand.Intn(len(pr_reviewers))
			pullRequest.AssignedReviwers[i] = pr_reviewers[idx]
			pr := s.pullRequestRepo.UpdatePR(pullRequest)
			return pr, nil
		}
	}

	return nil, entities.ErrNotAssigned
}
