package models

import (
    "time"
    "gorm.io/gorm"
)

type ChatMessage struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Content   string         `json:"content" gorm:"not null"`
    UserID    uint           `json:"user_id" gorm:"not null"`
    ProjectID uint           `json:"project_id" gorm:"not null"`
    FileID    *uint          `json:"file_id,omitempty"` // Optional, for file-specific messages
    TaskID    *uint          `json:"task_id,omitempty"` // Optional, for task-specific messages
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    User    User     `json:"user" gorm:"foreignKey:UserID"`
    Project Project  `json:"project" gorm:"foreignKey:ProjectID"`
    File    *File    `json:"file,omitempty" gorm:"foreignKey:FileID"`
    Task    *Task    `json:"task,omitempty" gorm:"foreignKey:TaskID"`
}