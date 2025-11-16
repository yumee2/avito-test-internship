package teams

import (
	"avito-test/internal/domain/entities"
	"avito-test/internal/domain/service/teams"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTeam struct {
	TeamName    string                `json:"team_name" binding:"required"`
	TeamMembers []entities.TeamMember `json:"members" binding:"required"`
}

type TeamController struct {
	teamService *teams.TeamService
}

func NewTeamController(teamService *teams.TeamService) *TeamController {
	return &TeamController{teamService: teamService}
}

func (c *TeamController) CreateTeam(ctx *gin.Context) {
	const fn = "adapters.http_server.team.CreateTeam"
	slog.With(slog.String("fn", fn))

	var newTeam createTeam
	if err := ctx.ShouldBindJSON(&newTeam); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	team, err := c.teamService.CreateTeam(newTeam.TeamName, newTeam.TeamMembers)
	if err != nil {
		if errors.Is(err, entities.ErrDuplicate) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    "TEAM_EXISTS",
				"message": "team_name already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, team)
}

func (c *TeamController) GetTeam(ctx *gin.Context) {
	teamName := ctx.Query("team_name")
	team, err := c.teamService.GetTeam(teamName)

	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    "NOT_FOUND",
				"message": "Team was not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, team)

}
