package models

import (
    "time"
    "gorm.io/gorm"
)

type TaskStatus string

const (
    TaskStatusTodo       TaskStatus = "todo"
    TaskStatusInProgress TaskStatus = "in_progress"
    TaskStatusDone       TaskStatus = "done"
)

type Task struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    ProjectID   uint           `json:"project_id"`
    SprintID    *uint          `json:"sprint_id"`
    AssigneeID  *uint          `json:"assignee_id"`
    Title       string         `json:"title" gorm:"not null"`
    Description string         `json:"description" gorm:"type:text"`
    Status      TaskStatus     `json:"status" gorm:"default:'todo'"`
    Priority    int            `json:"priority" gorm:"default:0"`
    GitHubIssue int            `json:"github_issue"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    Project  Project  `json:"project"`
    Sprint   *Sprint  `json:"sprint"`
    Assignee *User    `json:"assignee"`
    Comments []Comment `json:"comments"`
}

type Sprint struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    ProjectID   uint           `json:"project_id"`
    Name        string         `json:"name" gorm:"not null"`
    Description string         `json:"description"`
    StartDate   time.Time      `json:"start_date"`
    EndDate     time.Time      `json:"end_date"`
    IsActive    bool           `json:"is_active" gorm:"default:false"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    Project Project `json:"project"`
    Tasks   []Task  `json:"tasks"`
}

type Comment struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    TaskID    uint           `json:"task_id"`
    UserID    uint           `json:"user_id"`
    Content   string         `json:"content" gorm:"type:text;not null"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    Task Task `json:"task"`
    User User `json:"user"`
}