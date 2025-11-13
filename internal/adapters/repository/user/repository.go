package userRepository

import (
	"avito-test/internal/config"
	"avito-test/internal/domain/entities"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepositoryPG struct {
	db *gorm.DB
}

func New(cfg *config.Config) (*UserRepositoryPG, error) {
	const fn = "adapters.repository.user_repository.New"
	dsn := fmt.Sprintf("host=%s user=%s "+
		"password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresConnect.Host, cfg.PostgresConnect.User, cfg.PostgresConnect.Password, cfg.PostgresConnect.DatabaseName, cfg.PostgresConnect.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = db.AutoMigrate(&entities.User{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &UserRepositoryPG{db: db}, nil
}
