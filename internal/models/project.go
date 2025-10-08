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
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Users         []User          `json:"users" gorm:"many2many:user_projects;"`
    Files         []File          `json:"files"`
    Tasks         []Task          `json:"tasks"`
    Sprints       []Sprint        `json:"sprints"`
    Documentation []Documentation `json:"documentation"`
    ChatMessages  []ChatMessage   `json:"chat_messages"`
    Deployments   []Deployment    `json:"deployments"`
}

type File struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    ProjectID uint           `json:"project_id"`
    Name      string         `json:"name" gorm:"not null"`
    Path      string         `json:"path" gorm:"not null"`
    Content   string         `json:"content" gorm:"type:text"`
    Language  string         `json:"language"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    Project Project `json:"project"`
}