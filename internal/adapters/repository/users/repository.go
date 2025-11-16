package user

import (
	"avito-test/internal/domain/entities"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRepositoryPG struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*UserRepositoryPG, error) {
	const fn = "adapters.repository.user.New"

	err := db.AutoMigrate(&entities.User{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &UserRepositoryPG{db: db}, nil
}

func (r *UserRepositoryPG) SetIsActive(userId string, isActive bool) (*entities.User, error) {
	const fn = "adapters.repository.user.SetIsActive"

	result := r.db.Model(&entities.User{}).
		Where("user_id = ?", userId).
		Updates(map[string]interface{}{
			"is_active": isActive,
		})

	if result.Error != nil {
		return nil, fmt.Errorf("%s: %w", fn, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, entities.ErrNotFound
	}

	var user entities.User

	err := r.db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &user, nil
}

func (r *UserRepositoryPG) CreateUser(user *entities.User) (*entities.User, error) {
	result := r.db.Create(&user)
	var pgErr *pgconn.PgError

	if result.Error != nil {
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return nil, entities.ErrDuplicate
		}
	}

	return user, nil
}

func (r *UserRepositoryPG) GetUserFromTheSameTeam(userID string) ([]entities.User, error) {
	var userTeamName string
	r.db.Model(&entities.User{}).Select("team_name").Where("user_id = ?", userID).First(&userTeamName)

	if userTeamName == "" {
		return nil, entities.ErrNotFound
	}

	var users []entities.User
	r.db.Model(&entities.User{}).Select("user_id").Where("team_name = ?", userTeamName).Where("is_active = ?", true).Where("user_id <> ?", userID).Find(&users)

	return users, nil
}

func (r *UserRepositoryPG) GetUserForReassign(reviewers []entities.User, authorId string) ([]entities.User, error) {
	var userTeamName string
	r.db.Model(&entities.User{}).Select("team_name").Where("user_id = ?", authorId).First(&userTeamName)

	if userTeamName == "" {
		return nil, entities.ErrNotFound
	}

	var users []entities.User
	if len(reviewers) >= 2 {
		r.db.Model(&entities.User{}).Select("user_id").Where("team_name = ?", userTeamName).Where("is_active = ?", true).Where("user_id <> ?", reviewers[0].UserID).Where("user_id <> ?", authorId).Where("user_id <> ?", reviewers[1].UserID).Find(&users)
	} else {
		r.db.Model(&entities.User{}).Select("user_id").Where("team_name = ?", userTeamName).Where("is_active = ?", true).Where("user_id <> ?", reviewers[0].UserID).Where("user_id <> ?", authorId).Find(&users)
	}

	return users, nil
}
