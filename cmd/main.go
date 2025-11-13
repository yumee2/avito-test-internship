package main

import (
	userRepository "avito-test/internal/adapters/repository/user"
	"avito-test/internal/config"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.MustLoad()

	userRepo, err := userRepository.New(cfg)
	if err != nil {
		slog.Error("failed to setup database connection")
		os.Exit(1)
	}

	_ = userRepo
}
