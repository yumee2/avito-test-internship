package pullrequest

import (
	"avito-test/internal/domain/entities"
	pullrequest "avito-test/internal/domain/service/pull_request"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type prCreate struct {
	PullRequestId   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

type prMerge struct {
	PullRequestId string `json:"pull_request_id" binding:"required"`
}

type reassign struct {
	PullRequestId string `json:"pull_request_id" binding:"required"`
	OldUserId     string `json:"old_user_id" binding:"required"`
}

type PullRequestController struct {
	prService *pullrequest.PullRequestService
}

func NewPullRequestController(prService *pullrequest.PullRequestService) *PullRequestController {
	return &PullRequestController{prService: prService}
}

func (c *PullRequestController) CreatePullRequest(ctx *gin.Context) {
	var prCreate prCreate
	if err := ctx.ShouldBindJSON(&prCreate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	pullRequest, err := c.prService.CreatePR(prCreate.PullRequestId, prCreate.PullRequestName, prCreate.AuthorID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    "NOT_FOUND",
				"message": "author or a team not found",
			})
			return
		} else if errors.Is(err, entities.ErrDuplicate) {
			ctx.JSON(http.StatusConflict, gin.H{
				"code":    "PR_EXISTS",
				"message": "PR id already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, pullRequest.ToResponse())
}

func (c *PullRequestController) MergePullRequest(ctx *gin.Context) {
	var prMerge prMerge
	if err := ctx.ShouldBindJSON(&prMerge); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	pullRequest, err := c.prService.MergePr(prMerge.PullRequestId)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    "NOT_FOUND",
				"message": "pr was not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, pullRequest.ToResponse())
}

func (c *PullRequestController) Reassing(ctx *gin.Context) {
	var reassignData reassign
	if err := ctx.ShouldBindJSON(&reassignData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	pullRequest, err := c.prService.Reassign(reassignData.PullRequestId, reassignData.OldUserId)
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    "NOT_FOUND",
				"message": "PR or user was not found",
			})
		case errors.Is(err, entities.ErrPrMergerd):
			ctx.JSON(http.StatusConflict, gin.H{
				"code":    "PR_MERGED",
				"message": "Cannot reassign on merged PR",
			})
		case errors.Is(err, entities.ErrNotAssigned):
			ctx.JSON(http.StatusConflict, gin.H{
				"code":    "NOT_ASSIGNED",
				"message": "Reviewer is not assigned to this PR",
			})
		case errors.Is(err, entities.ErrNoCandidates):
			ctx.JSON(http.StatusConflict, gin.H{
				"code":    "NO_CANDIDATE",
				"message": "No active replacement candidate in team",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, pullRequest.ToResponse())
}
