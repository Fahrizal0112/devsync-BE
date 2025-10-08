package handlers

import (
	"net/http"
	"strconv"

	"devsync-be/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	db *gorm.DB
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{db: db}
}

// @Summary Get projects
// @Description Get all projects for authenticated user
// @Tags projects
// @Security BearerAuth
// @Success 200 {array} models.Project
// @Router /projects [get]
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var projects []models.Project
	err := h.db.Joins("JOIN user_projects ON user_projects.project_id = projects.id").
		Where("user_projects.user_id = ?", userID).
		Preload("Users").
		Find(&projects).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// @Summary Create project
// @Description Create a new project
// @Tags projects
// @Security BearerAuth
// @Param project body models.Project true "Project data"
// @Success 201 {object} models.Project
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create project
	if err := h.db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Add user to project
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	if err := h.db.Model(&project).Association("Users").Append(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// @Summary Get project
// @Description Get project by ID
// @Tags projects
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	if err := h.db.Preload("Users").Preload("Files").Preload("Tasks").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// @Summary Update project
// @Description Update project by ID
// @Tags projects
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param project body models.Project true "Project data"
// @Success 200 {object} models.Project
// @Router /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	if err := h.db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// @Summary Delete project
// @Description Delete project by ID
// @Tags projects
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 204
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.db.Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.Status(http.StatusNoContent)
}
