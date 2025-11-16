package teams

import "avito-test/internal/domain/entities"

type teamRepository interface {
	CreateTeam(team *entities.Team) (*entities.Team, error)
	GetTeam(teamName string) (*entities.Team, error)
}

type userRepository interface {
	CreateUser(user *entities.User) (*entities.User, error)
}
