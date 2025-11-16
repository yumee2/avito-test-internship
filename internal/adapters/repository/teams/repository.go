package teams

import (
	"avito-test/internal/domain/entities"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type TeamRepositoryPG struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*TeamRepositoryPG, error) {
	const fn = "adapters.repository.team.New"

	err := db.AutoMigrate(&entities.Team{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &TeamRepositoryPG{db: db}, nil
}

func (t *TeamRepositoryPG) CreateTeam(team *entities.Team) (*entities.Team, error) {
	const fn = "adapters.repository.teams.CreateTeam"

	result := t.db.Create(&team)
	var pgErr *pgconn.PgError

	if result.Error != nil {
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return nil, entities.ErrDuplicate
		}
		return nil, fmt.Errorf("%s: %w", fn, result.Error)
	}

	return team, nil
}

func (t *TeamRepositoryPG) GetTeam(teamName string) (*entities.Team, error) {
	const fn = "adapters.repository.teams.GetTeam"
	var team entities.Team

	err := t.db.
		Preload("TeamMembers").
		Where("team_name = ?", teamName).
		First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}
