package users

import "github.com/gin-gonic/gin"

func SetupUsersRoutes(r *gin.Engine, userController *UserController) {
	r.POST("/users/setIsActive", userController.SetIsActive)
	r.GET("/users/getReview", userController.GetPullRequestsForUser)
}
