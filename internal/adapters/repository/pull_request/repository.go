package pull_request

import (
	"avito-test/internal/domain/entities"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type PRRepositoryPG struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*PRRepositoryPG, error) {
	const fn = "adapters.repository.pull_request.New"

	err := db.AutoMigrate(&entities.PullRequest{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &PRRepositoryPG{db: db}, nil
}

func (p *PRRepositoryPG) CreatePR(pr *entities.PullRequest) (*entities.PullRequest, error) {
	const fn = "adapters.repository.pull_request.CreatePR"

	result := p.db.Create(&pr)
	var pgErr *pgconn.PgError

	if result.Error != nil {
		if errors.As(result.Error, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return nil, entities.ErrDuplicate
			case "23503": // foreign_key_violation
				return nil, entities.ErrNotFound
			}
		}
		return nil, fmt.Errorf("%s: %w", fn, result.Error)
	}

	return pr, nil
}

func (p *PRRepositoryPG) MergePr(pullRequestID string) (*entities.PullRequest, error) {
	const fn = "adapters.repository.pull_request.MergePr"

	var existingPr entities.PullRequest

	err := p.db.Where("pull_request_id = ?", pullRequestID).First(&existingPr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if existingPr.Status == string(entities.StatusMerged) {
		err := p.db.Preload("AssignedReviwers").First(&existingPr).Error
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		return &existingPr, nil
	}

	mergedTime := time.Now()
	result := p.db.Model(&entities.PullRequest{}).
		Where("pull_request_id = ?", pullRequestID).
		Updates(entities.PullRequest{Status: string(entities.StatusMerged), MergedAt: &mergedTime})

	if result.Error != nil {
		return nil, fmt.Errorf("%s: %w", fn, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, entities.ErrNotFound
	}

	var pr entities.PullRequest

	err = p.db.Preload("AssignedReviwers").Where("pull_request_id = ?", pullRequestID).First(&pr).Error
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &pr, nil

}

func (p *PRRepositoryPG) GetPR(pullRequestID string) (*entities.PullRequest, error) {
	var existingPr entities.PullRequest

	err := p.db.Preload("AssignedReviwers").Where("pull_request_id = ?", pullRequestID).First(&existingPr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}

	return &existingPr, nil
}

func (p *PRRepositoryPG) UpdatePR(pr *entities.PullRequest) *entities.PullRequest {
	p.db.Model(&pr).
		Association("AssignedReviwers").
		Replace(pr.AssignedReviwers)

	var pullRequest entities.PullRequest

	p.db.Preload("AssignedReviwers").Where("pull_request_id = ?", pr.PullRequestID).First(&pullRequest)
	return &pullRequest
}

func (p *PRRepositoryPG) GetPullRequestsForUser(userID string) (*[]entities.PullRequest, error) {
	var pullRequests *[]entities.PullRequest

	err := p.db.
		Preload("AssignedReviwers").
		Joins("LEFT JOIN pr_reviewers ON pr_reviewers.pull_request_pull_request_id = pull_requests.pull_request_id").
		Where("pr_reviewers.user_user_id = ?", userID).
		Find(&pullRequests).Error

	if err != nil {
		return nil, err
	}

	return pullRequests, err
}
