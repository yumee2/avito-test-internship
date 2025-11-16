package pullrequest

import "avito-test/internal/domain/entities"

type pullRequestRepo interface {
	CreatePR(pr *entities.PullRequest) (*entities.PullRequest, error)
	MergePr(pullRequestID string) (*entities.PullRequest, error)
	GetPR(pullRequestID string) (*entities.PullRequest, error)
	UpdatePR(pr *entities.PullRequest) *entities.PullRequest
}

type userRepository interface {
	GetUserFromTheSameTeam(userID string) ([]entities.User, error)
	GetUserForReassign(userIDs []entities.User, authorId string) ([]entities.User, error)
}
