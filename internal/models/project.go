package models

import (
    "time"
    "gorm.io/gorm"
)

type Project struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"not null"`
    Description string         `json:"description"`
    GitHubRepo  string         `json:"github_repo"`
    IsPublic    bool           `json:"is_public" gorm:"default:false"`
    CreatedBy   *uint          `json:"created_by"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Creator       User            `json:"creator" gorm:"foreignKey:CreatedBy"`
    Users         []User          `json:"users" gorm:"many2many:user_projects;"`
    Files         []File          `json:"files" gorm:"foreignKey:ProjectID"`
    Tasks         []Task          `json:"tasks" gorm:"foreignKey:ProjectID"`
    Sprints       []Sprint        `json:"sprints" gorm:"foreignKey:ProjectID"`
    Documentation []Documentation `json:"documentation" gorm:"foreignKey:ProjectID"`
    ChatMessages  []ChatMessage   `json:"chat_messages" gorm:"foreignKey:ProjectID"`
    Deployments   []Deployment    `json:"deployments" gorm:"foreignKey:ProjectID"`
}