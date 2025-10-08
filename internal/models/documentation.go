package models

import (
    "time"
    "gorm.io/gorm"
)

type Documentation struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    ProjectID uint           `json:"project_id"`
    Title     string         `json:"title" gorm:"not null"`
    Content   string         `json:"content" gorm:"type:text"`
    Path      string         `json:"path"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    Project Project `json:"project" gorm:"foreignKey:ProjectID"`
}