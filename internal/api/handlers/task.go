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

type TaskHandler struct {
    db  *gorm.DB
    hub *websocket.Hub
}

func NewTaskHandler(db *gorm.DB, hub *websocket.Hub) *TaskHandler {
    return &TaskHandler{
        db:  db,
        hub: hub,
    }
}

// @Summary Get tasks
// @Description Get all tasks in a project
// @Tags tasks
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {array} models.Task
// @Router /projects/{id}/tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    var tasks []models.Task
    if err := h.db.Where("project_id = ?", projectID).
        Preload("Assignee").
        Preload("Sprint").
        Preload("Comments").
        Find(&tasks).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
        return
    }

    c.JSON(http.StatusOK, tasks)
}

// @Summary Create task
// @Description Create a new task in project
// @Tags tasks
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param task body models.Task true "Task data"
// @Success 201 {object} models.Task
// @Router /projects/{id}/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    var task models.Task
    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    task.ProjectID = uint(projectID)

    if err := h.db.Create(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }

    // Load relationships
    h.db.Preload("Assignee").Preload("Sprint").First(&task, task.ID)

    // Broadcast task creation to WebSocket clients
    message := map[string]interface{}{
        "type":       "task_created",
        "project_id": projectID,
        "data":       task,
    }
    if msgBytes, err := json.Marshal(message); err == nil {
        h.hub.Broadcast(msgBytes)
    }

    c.JSON(http.StatusCreated, task)
}

// @Summary Update task
// @Description Update task by ID
// @Tags tasks
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param taskId path int true "Task ID"
// @Param task body models.Task true "Task data"
// @Success 200 {object} models.Task
// @Router /projects/{id}/tasks/{taskId} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    taskID, err := strconv.Atoi(c.Param("taskId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    var task models.Task
    if err := h.db.First(&task, taskID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.db.Save(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
        return
    }

    // Load relationships
    h.db.Preload("Assignee").Preload("Sprint").First(&task, task.ID)

    // Broadcast task update to WebSocket clients
    message := map[string]interface{}{
        "type":       "task_updated",
        "project_id": projectID,
        "data":       task,
    }
    if msgBytes, err := json.Marshal(message); err == nil {
        h.hub.Broadcast(msgBytes)
    }

    c.JSON(http.StatusOK, task)
}

// @Summary Delete task
// @Description Delete task by ID
// @Tags tasks
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param taskId path int true "Task ID"
// @Success 204
// @Router /projects/{id}/tasks/{taskId} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    taskID, err := strconv.Atoi(c.Param("taskId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    if err := h.db.Delete(&models.Task{}, taskID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
        return
    }

    // Broadcast task deletion to WebSocket clients
    message := map[string]interface{}{
        "type":       "task_deleted",
        "project_id": projectID,
        "data":       map[string]interface{}{"id": taskID},
    }
    if msgBytes, err := json.Marshal(message); err == nil {
        h.hub.Broadcast(msgBytes)
    }

    c.Status(http.StatusNoContent)
}

// @Summary Get sprints
// @Description Get all sprints in a project
// @Tags sprints
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {array} models.Sprint
// @Router /projects/{id}/sprints [get]
func (h *TaskHandler) GetSprints(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    var sprints []models.Sprint
    if err := h.db.Where("project_id = ?", projectID).
        Preload("Tasks").
        Find(&sprints).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sprints"})
        return
    }

    c.JSON(http.StatusOK, sprints)
}

// @Summary Create sprint
// @Description Create a new sprint in project
// @Tags sprints
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param sprint body models.Sprint true "Sprint data"
// @Success 201 {object} models.Sprint
// @Router /projects/{id}/sprints [post]
func (h *TaskHandler) CreateSprint(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    var sprint models.Sprint
    if err := c.ShouldBindJSON(&sprint); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    sprint.ProjectID = uint(projectID)

    if err := h.db.Create(&sprint).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sprint"})
        return
    }

    c.JSON(http.StatusCreated, sprint)
}