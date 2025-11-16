package main

import (
	pullrequestHttp "avito-test/internal/adapters/http_server/pull_request"
	teamHttp "avito-test/internal/adapters/http_server/teams"
	userHttp "avito-test/internal/adapters/http_server/users"
	pullRequestRepo "avito-test/internal/adapters/repository/pull_request"
	teamRepository "avito-test/internal/adapters/repository/teams"
	userRepository "avito-test/internal/adapters/repository/users"
	"avito-test/internal/config"
	pullrequestService "avito-test/internal/domain/service/pull_request"
	teamService "avito-test/internal/domain/service/teams"
	usersService "avito-test/internal/domain/service/users"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.MustLoad()
	db := mustInitDB(cfg)

	userRepo, err := userRepository.New(db)
	if err != nil {
		slog.Error("failed to init a user repository")
		os.Exit(1)
	}

	prRepo, err := pullRequestRepo.New(db)
	if err != nil {
		slog.Error("failed to init a pull request repository")
		os.Exit(1)
	}

	teamRepo, err := teamRepository.New(db)
	if err != nil {
		slog.Error("failed to init a pull request repository")
		os.Exit(1)
	}

	teamService := teamService.NewTeamService(teamRepo, userRepo)
	userService := usersService.NewUserService(userRepo, prRepo)
	pullRequestService := pullrequestService.NewPullRequestService(prRepo, userRepo)

	r := setUpHttpServer(teamService, userService, pullRequestService)

	if err := r.Run(cfg.Address); err != nil {
		slog.Error("Failed to start server:", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
	}
}

func setUpHttpServer(teamService *teamService.TeamService, userService *usersService.UserService, prService *pullrequestService.PullRequestService) *gin.Engine {
	r := gin.Default()

	teamController := teamHttp.NewTeamController(teamService)
	teamHttp.SetupTeamsRoutes(r, teamController)

	userController := userHttp.NewUserController(userService)
	userHttp.SetupUsersRoutes(r, userController)

	pullRequestController := pullrequestHttp.NewPullRequestController(prService)
	pullrequestHttp.SetupPRRoutes(r, pullRequestController)

	return r
}

func mustInitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s "+
		"password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresConnect.Host, cfg.PostgresConnect.User, cfg.PostgresConnect.Password, cfg.PostgresConnect.DatabaseName, cfg.PostgresConnect.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		slog.Error("failed to open a database connection")
		os.Exit(1)
	}
	return db
}
