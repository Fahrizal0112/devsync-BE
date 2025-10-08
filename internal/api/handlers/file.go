package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "devsync-be/internal/models"
    "devsync-be/internal/websocket"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type FileHandler struct {
    db  *gorm.DB
    hub *websocket.Hub
}

func NewFileHandler(db *gorm.DB, hub *websocket.Hub) *FileHandler {
    return &FileHandler{
        db:  db,
        hub: hub,
    }
}

// @Summary Get files
// @Description Get all files in a project
// @Tags files
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {array} models.File
// @Router /projects/{id}/files [get]
func (h *FileHandler) GetFiles(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    var files []models.File
    if err := h.db.Where("project_id = ?", projectID).Find(&files).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
        return
    }

    c.JSON(http.StatusOK, files)
}

// @Summary Create file
// @Description Create a new file in project
// @Tags files
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param file body models.File true "File data"
// @Success 201 {object} models.File
// @Router /projects/{id}/files [post]
func (h *FileHandler) CreateFile(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    var file models.File
    if err := c.ShouldBindJSON(&file); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Validate required fields
    if file.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File name is required"})
        return
    }
    if file.Path == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
        return
    }

    // Set project ID and uploaded_by from context
    file.ProjectID = uint(projectID)
    
    // Get user ID from JWT token context
    if userID, exists := c.Get("userID"); exists {
        file.UploadedBy = userID.(uint)
    }

    // Verify project exists
    var project models.Project
    if err := h.db.First(&project, projectID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
        return
    }

    if err := h.db.Create(&file).Error; err != nil {
        // Log the actual error for debugging
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create file",
            "details": err.Error(),
        })
        return
    }

    // Broadcast file creation to WebSocket clients
    message := map[string]interface{}{
        "type":       "file_created",
        "project_id": projectID,
        "data":       file,
    }
    if msgBytes, err := json.Marshal(message); err == nil {
        h.hub.Broadcast(msgBytes)
    }

    c.JSON(http.StatusCreated, file)
}

// @Summary Get file
// @Description Get file by ID
// @Tags files
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param fileId path int true "File ID"
// @Success 200 {object} models.File
// @Router /projects/{id}/files/{fileId} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
    fileID, err := strconv.Atoi(c.Param("fileId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
        return
    }

    var file models.File
    if err := h.db.First(&file, fileID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
        return
    }

    c.JSON(http.StatusOK, file)
}

// @Summary Update file
// @Description Update file content
// @Tags files
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param fileId path int true "File ID"
// @Param file body models.File true "File data"
// @Success 200 {object} models.File
// @Router /projects/{id}/files/{fileId} [put]
func (h *FileHandler) UpdateFile(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    fileID, err := strconv.Atoi(c.Param("fileId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
        return
    }

    var file models.File
    if err := h.db.First(&file, fileID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
        return
    }

    if err := c.ShouldBindJSON(&file); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.db.Save(&file).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file"})
        return
    }

    // Broadcast file update to WebSocket clients
    message := map[string]interface{}{
        "type":       "file_updated",
        "project_id": projectID,
        "data":       file,
    }
    if msgBytes, err := json.Marshal(message); err == nil {
        h.hub.Broadcast(msgBytes)
    }

    c.JSON(http.StatusOK, file)
}

// @Summary Delete file
// @Description Delete file by ID
// @Tags files
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param fileId path int true "File ID"
// @Success 204
// @Router /projects/{id}/files/{fileId} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    fileID, err := strconv.Atoi(c.Param("fileId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
        return
    }

    if err := h.db.Delete(&models.File{}, fileID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
        return
    }

    // Broadcast file deletion to WebSocket clients
    message := map[string]interface{}{
        "type":       "file_deleted",
        "project_id": projectID,
        "data":       map[string]interface{}{"id": fileID},
    }
    if msgBytes, err := json.Marshal(message); err == nil {
        h.hub.Broadcast(msgBytes)
    }

    c.Status(http.StatusNoContent)
}