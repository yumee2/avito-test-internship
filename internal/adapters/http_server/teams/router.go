package teams

import "github.com/gin-gonic/gin"

func SetupTeamsRoutes(r *gin.Engine, teamController *TeamController) {
	r.POST("/team/add", teamController.CreateTeam)
	r.GET("/team/get", teamController.GetTeam)
}
