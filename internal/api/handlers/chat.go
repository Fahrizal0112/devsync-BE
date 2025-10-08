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

type ChatHandler struct {
    db  *gorm.DB
    hub *websocket.Hub
}

func NewChatHandler(db *gorm.DB, hub *websocket.Hub) *ChatHandler {
    return &ChatHandler{
        db:  db,
        hub: hub,
    }
}

// @Summary Get messages
// @Description Get all chat messages in a project
// @Tags chat
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param file_id query int false "Filter by file ID"
// @Param task_id query int false "Filter by task ID"
// @Success 200 {array} models.ChatMessage
// @Router /projects/{id}/messages [get]
func (h *ChatHandler) GetMessages(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    query := h.db.Where("project_id = ?", projectID).
        Preload("User").
        Preload("File").
        Preload("Task")

    // Filter by file if provided
    if fileID := c.Query("file_id"); fileID != "" {
        query = query.Where("file_id = ?", fileID)
    }

    // Filter by task if provided
    if taskID := c.Query("task_id"); taskID != "" {
        query = query.Where("task_id = ?", taskID)
    }

    var messages []models.ChatMessage
    if err := query.Order("created_at ASC").Find(&messages).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
        return
    }

    c.JSON(http.StatusOK, messages)
}

// @Summary Send message
// @Description Send a chat message in project
// @Tags chat
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param message body models.ChatMessage true "Message data"
// @Success 201 {object} models.ChatMessage
// @Router /projects/{id}/messages [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    var message models.ChatMessage
    if err := c.ShouldBindJSON(&message); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    message.ProjectID = uint(projectID)
    message.UserID = userID.(uint)

    if err := h.db.Create(&message).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
        return
    }

    // Load relationships
    h.db.Preload("User").Preload("File").Preload("Task").First(&message, message.ID)

    // Broadcast message to WebSocket clients
    wsMessage := map[string]interface{}{
        "type":       "chat_message",
        "project_id": projectID,
        "data":       message,
    }
    if msgBytes, err := json.Marshal(wsMessage); err == nil {
        h.hub.Broadcast(msgBytes)
    }

    c.JSON(http.StatusCreated, message)
}