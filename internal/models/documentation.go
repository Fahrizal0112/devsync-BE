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

    Project Project `json:"project"`
}

type ChatMessage struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    ProjectID uint           `json:"project_id"`
    UserID    uint           `json:"user_id"`
    FileID    *uint          `json:"file_id"`
    TaskID    *uint          `json:"task_id"`
    Content   string         `json:"content" gorm:"type:text;not null"`
    ThreadID  *uint          `json:"thread_id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    Project Project `json:"project"`
    User    User    `json:"user"`
    File    *File   `json:"file"`
    Task    *Task   `json:"task"`
}

type Deployment struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    ProjectID   uint           `json:"project_id"`
    Branch      string         `json:"branch" gorm:"not null"`
    CommitHash  string         `json:"commit_hash"`
    Status      string         `json:"status" gorm:"default:'pending'"`
    PreviewURL  string         `json:"preview_url"`
    LogsURL     string         `json:"logs_url"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    Project Project `json:"project"`
}