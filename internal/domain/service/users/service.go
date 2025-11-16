package users

import (
	"avito-test/internal/domain/entities"
	"errors"
)

type UserService struct {
	userRepo        userRepository
	pullRequestRepo prRepo
}

func NewUserService(userRepo userRepository, pullRequestRepo prRepo) *UserService {
	return &UserService{userRepo: userRepo, pullRequestRepo: pullRequestRepo}
}

func (s *UserService) SetIsActive(userId string, isActive bool) (*entities.User, error) {
	user, err := s.userRepo.SetIsActive(userId, isActive)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetPullRequestsForUser(userID string) (*[]entities.PullRequestShort, error) {
	var usersPullRequests []entities.PullRequestShort
	pullRequests, err := s.pullRequestRepo.GetPullRequestsForUser(userID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.ErrNotFound
		}

		return nil, err
	}

	for _, v := range *pullRequests {
		usersPullRequests = append(usersPullRequests, v.ToShortVersion())
	}

	return &usersPullRequests, nil
}
