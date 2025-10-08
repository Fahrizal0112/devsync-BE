package models

import (
    "time"
    "gorm.io/gorm"
)

type Deployment struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Version     string         `json:"version" gorm:"not null"`
    Status      string         `json:"status" gorm:"not null"` // pending, deploying, success, failed
    Environment string         `json:"environment" gorm:"not null"` // development, staging, production
    URL         string         `json:"url"`
    ProjectID   uint           `json:"project_id" gorm:"not null"`
    DeployedBy  uint           `json:"deployed_by" gorm:"not null"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Project    Project `json:"project" gorm:"foreignKey:ProjectID"`
    DeployedByUser User `json:"deployed_by_user" gorm:"foreignKey:DeployedBy"`
}