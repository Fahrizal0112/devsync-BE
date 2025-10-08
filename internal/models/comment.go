package models

import (
    "time"
    "gorm.io/gorm"
)

type Comment struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Content   string         `json:"content" gorm:"not null"`
    UserID    uint           `json:"user_id" gorm:"not null"`
    TaskID    uint           `json:"task_id" gorm:"not null"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    User User `json:"user" gorm:"foreignKey:UserID"`
    Task Task `json:"task" gorm:"foreignKey:TaskID"`
}