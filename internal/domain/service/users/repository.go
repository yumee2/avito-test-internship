package users

import "avito-test/internal/domain/entities"

type userRepository interface {
	SetIsActive(userId string, isActive bool) (*entities.User, error)
}

type prRepo interface {
	GetPullRequestsForUser(userID string) (*[]entities.PullRequest, error)
}
