package models

import (
    "time"
    "gorm.io/gorm"
)

type Sprint struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"not null"`
    Description string         `json:"description"`
    StartDate   time.Time      `json:"start_date"`
    EndDate     time.Time      `json:"end_date"`
    Status      string         `json:"status" gorm:"default:'active'"` // active, completed, cancelled
    ProjectID   uint           `json:"project_id" gorm:"not null"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Project Project `json:"project" gorm:"foreignKey:ProjectID"`
    Tasks   []Task  `json:"tasks" gorm:"foreignKey:SprintID"`
}