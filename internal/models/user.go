package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    GitHubID    int64          `json:"github_id" gorm:"uniqueIndex"`
    Username    string         `json:"username" gorm:"uniqueIndex;not null"`
    Email       string         `json:"email" gorm:"uniqueIndex;not null"`
    Name        string         `json:"name"`
    AvatarURL   string         `json:"avatar_url"`
    AccessToken string         `json:"-"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Projects     []Project     `json:"projects" gorm:"many2many:user_projects;"`
    Tasks        []Task        `json:"tasks" gorm:"foreignKey:AssigneeID"`
    Comments     []Comment     `json:"comments" gorm:"foreignKey:UserID"`
    ChatMessages []ChatMessage `json:"chat_messages" gorm:"foreignKey:UserID"`
}