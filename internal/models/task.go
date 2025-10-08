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

    // Relationships
    Project  Project  `json:"project" gorm:"foreignKey:ProjectID"`
    Sprint   *Sprint  `json:"sprint" gorm:"foreignKey:SprintID"`
    Assignee *User    `json:"assignee" gorm:"foreignKey:AssigneeID"`
    Comments []Comment `json:"comments" gorm:"foreignKey:TaskID"`
}