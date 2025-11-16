package pullrequest

import "github.com/gin-gonic/gin"

func SetupPRRoutes(r *gin.Engine, pullRequestController *PullRequestController) {
	r.POST("/pullRequest/create", pullRequestController.CreatePullRequest)
	r.POST("/pullRequest/merge", pullRequestController.MergePullRequest)
	r.POST("/pullRequest/reassign", pullRequestController.Reassing)
}
