package database

import (
	"devsync-be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

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