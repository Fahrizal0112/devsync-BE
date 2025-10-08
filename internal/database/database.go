package database

import (
	"devsync-be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.File{},
		&models.Task{},
		&models.Sprint{},
		&models.Comment{},
		&models.Documentation{},
		&models.ChatMessage{},
		&models.Deployment{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}