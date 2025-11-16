package users

import (
	"avito-test/internal/domain/entities"
	"avito-test/internal/domain/service/users"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userChangeState struct {
	UserId   string `json:"user_id" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type pullReviewsResponse struct {
	UserID       string                      `json:"user_id" binding:"required"`
	PullRequests []entities.PullRequestShort `json:"pull_requests" binding:"required"`
}

type UserController struct {
	userService *users.UserService
}

func NewUserController(userService *users.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) SetIsActive(ctx *gin.Context) {
	var userChangeState userChangeState
	if err := ctx.ShouldBindJSON(&userChangeState); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	slog.Info(userChangeState.UserId)
	user, err := c.userService.SetIsActive(userChangeState.UserId, *userChangeState.IsActive)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    "NOT_FOUND",
				"message": "user not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) GetPullRequestsForUser(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	pullRequests, err := c.userService.GetPullRequestsForUser(userID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    "NOT_FOUND",
				"message": "user not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, pullReviewsResponse{UserID: userID, PullRequests: *pullRequests})
}
