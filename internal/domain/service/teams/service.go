package teams

import (
	"avito-test/internal/domain/entities"
	"errors"
)

type TeamService struct {
	teamRepo teamRepository
	userRepo userRepository
}

func NewTeamService(teamRepo teamRepository, userRepo userRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo}
}

func (s *TeamService) CreateTeam(teamName string, teamMembers []entities.TeamMember) (*entities.Team, error) {
	var users []entities.User
	for _, val := range teamMembers {
		user := entities.User{
			UserID:   val.UserId,
			Username: val.Username,
			TeamName: teamName,
			IsActive: val.IsActive,
		}
		createdUser, err := s.userRepo.CreateUser(&user)
		if err != nil {
			continue
		}
		users = append(users, *createdUser)
	}

	team := &entities.Team{
		TeamName:    teamName,
		TeamMembers: users,
	}

	team, err := s.teamRepo.CreateTeam(team)
	if err != nil {
		if errors.Is(err, entities.ErrDuplicate) {
			return nil, entities.ErrDuplicate
		}
		return nil, err
	}

	return team, nil
}

func (s *TeamService) GetTeam(teamName string) (*entities.Team, error) {
	team, err := s.teamRepo.GetTeam(teamName)
	if err != nil {
		return nil, entities.ErrNotFound
	}

	return team, nil
}
