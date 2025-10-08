package models

import (
    "time"
    "gorm.io/gorm"
)

type File struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"not null"`
    Path      string         `json:"path" gorm:"not null"`
    Content   string         `json:"content" gorm:"type:text"`
    FileURL   string         `json:"file_url"`
    FileType  string         `json:"file_type"`
    FileSize  int64          `json:"file_size"`
    MimeType  string         `json:"mime_type"`
    ProjectID uint           `json:"project_id" gorm:"not null"`
    UploadedBy uint          `json:"uploaded_by"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Project      Project       `json:"project" gorm:"foreignKey:ProjectID"`
    ChatMessages []ChatMessage `json:"chat_messages" gorm:"foreignKey:FileID"`
    Uploader     User          `json:"uploader" gorm:"foreignKey:UploadedBy"`
}