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

// Helper function to check if user is a member of the project
func (h *ProjectHandler) isProjectMember(userID, projectID uint) bool {
	var count int64
	h.db.Table("user_projects").
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Count(&count)
	return count > 0
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

	// Set creator
	creatorID := userID.(uint)
	project.CreatedBy = &creatorID

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
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check if user is a member of the project
	if !h.isProjectMember(userID.(uint), uint(id)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You are not a member of this project"})
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
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check if user is a member of the project
	if !h.isProjectMember(userID.(uint), uint(id)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You are not a member of this project"})
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
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check if project exists and user is the creator
	var project models.Project
	if err := h.db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Only creator can delete the project
	if project.CreatedBy == nil || *project.CreatedBy != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only project creator can delete this project"})
		return
	}

	if err := h.db.Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get project members
// @Description Get all members of a project
// @Tags projects
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {array} models.User
// @Router /projects/{id}/members [get]
func (h *ProjectHandler) GetMembers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check if user is a member of the project
	if !h.isProjectMember(userID.(uint), uint(projectID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. You are not a member of this project"})
		return
	}

	// Check if project exists
	var project models.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Get all members of the project
	var users []models.User
	err = h.db.Joins("JOIN user_projects ON user_projects.user_id = users.id").
		Where("user_projects.project_id = ?", projectID).
		Select("users.id, users.username, users.name, users.email, users.avatar_url, users.created_at").
		Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project members"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// AddMemberRequest represents the request body for adding a member
type AddMemberRequest struct {
	UserID   *uint  `json:"user_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

// @Summary Add member to project
// @Description Add a user to project by user_id, email, or username
// @Tags projects
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param member body AddMemberRequest true "Member data"
// @Success 201 {object} models.User
// @Router /projects/{id}/members [post]
func (h *ProjectHandler) AddMember(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check if user is a member of the project
	if !h.isProjectMember(userID.(uint), uint(projectID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. You are not a member of this project"})
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if project exists
	var project models.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Find user by different criteria
	var user models.User
	var query *gorm.DB

	if req.UserID != nil {
		query = h.db.Where("id = ?", *req.UserID)
	} else if req.Email != "" {
		query = h.db.Where("email = ?", req.Email)
	} else if req.Username != "" {
		query = h.db.Where("username = ?", req.Username)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide user_id, email, or username"})
		return
	}

	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user is already a member
	var count int64
	h.db.Table("user_projects").
		Where("user_id = ? AND project_id = ?", user.ID, projectID).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already a member of this project"})
		return
	}

	// Add user to project
	if err := h.db.Model(&project).Association("Users").Append(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to project"})
		return
	}

	// Return user info (without sensitive data)
	userResponse := models.User{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully added to project",
		"user":    userResponse,
	})
}

// @Summary Remove member from project
// @Description Remove a user from project
// @Tags projects
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param userId path int true "User ID"
// @Success 204
// @Router /projects/{id}/members/{userId} [delete]
func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"}) 
		return
	}

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if project exists
	var project models.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Only creator can remove members
	if project.CreatedBy == nil || *project.CreatedBy != currentUserID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only project creator can remove members"})
		return
	}

	// Creator cannot remove themselves
	if project.CreatedBy != nil && uint(userID) == *project.CreatedBy {
		c.JSON(http.StatusForbidden, gin.H{"error": "Project creator cannot be removed from the project"})
		return
	}

	// Check if user exists
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user is a member of the project
	var count int64
	h.db.Table("user_projects").
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Count(&count)

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User is not a member of this project"})
		return
	}

	// Remove user from project
	if err := h.db.Model(&project).Association("Users").Delete(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from project"})
		return
	}

	c.Status(http.StatusNoContent)
}
