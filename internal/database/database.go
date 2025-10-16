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

	// Migrate existing projects without CreatedBy
	err = migrateExistingProjects(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateExistingProjects(db *gorm.DB) error {
	// Find projects without CreatedBy
	var projects []models.Project
	err := db.Where("created_by IS NULL").Find(&projects).Error
	if err != nil {
		return err
	}

	// For each project, set the first user as creator
	for _, project := range projects {
		var firstUser models.User
		err := db.Joins("JOIN user_projects ON user_projects.user_id = users.id").
			Where("user_projects.project_id = ?", project.ID).
			Order("user_projects.created_at ASC").
			First(&firstUser).Error
		
		if err == nil {
			// Set the first user as creator
			project.CreatedBy = &firstUser.ID
			db.Save(&project)
		}
	}

	return nil
}