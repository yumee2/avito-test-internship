package entities

import (
	"time"

	"gorm.io/gorm"
)

type PullRequestStatus string

const (
	StatusOpen   PullRequestStatus = "OPEN"
	StatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID    string     `gorm:"primaryKey;default:gen_random_uuid()"`
	PullRequestName  string     `gorm:"size:16;not null;"`
	AuthorID         string     `gorm:"type:uuid;not null"`
	AssignedReviwers []User     `gorm:"many2many:pr_reviewers;"`
	Status           string     `gorm:"size:255;not null"`
	Author           User       `gorm:"foreignKey:AuthorID;references:UserID" json:"-"`
	CreatedAt        *time.Time `gorm:"column:created_at"`
	MergedAt         *time.Time `gorm:"column:merged_at"`
}

type PullRequestResponse struct {
	PullRequestID       string     `json:"pull_request_id"`
	PullRequestName     string     `json:"pull_request_name"`
	AuthorID            string     `json:"author_id"`
	AssignedReviewerIDs []string   `json:"assigned_reviewers"` // Только ID ревьюверов
	Status              string     `json:"status"`
	CreatedAt           *time.Time `json:"createdAt"`
	MergedAt            *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func (pr *PullRequest) ToResponse() PullRequestResponse {
	reviewerIDs := make([]string, len(pr.AssignedReviwers))
	for i, reviewer := range pr.AssignedReviwers {
		reviewerIDs[i] = reviewer.UserID
	}

	return PullRequestResponse{
		PullRequestID:       pr.PullRequestID,
		PullRequestName:     pr.PullRequestName,
		AuthorID:            pr.AuthorID,
		AssignedReviewerIDs: reviewerIDs,
		Status:              pr.Status,
		CreatedAt:           pr.CreatedAt,
		MergedAt:            pr.MergedAt,
	}
}

func (pr *PullRequest) ToShortVersion() PullRequestShort {
	return PullRequestShort{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          pr.Status,
	}
}

func (pr *PullRequest) BeforeCreate(tx *gorm.DB) error {
	if pr.CreatedAt == nil {
		now := time.Now()
		pr.CreatedAt = &now
	}
	return nil
}
