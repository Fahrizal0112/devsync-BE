package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"devsync-be/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// @Summary Search users
// @Description Search users by username or email
// @Tags users
// @Security BearerAuth
// @Param q query string true "Search query (username or email)"
// @Param limit query int false "Limit results (default: 10, max: 50)"
// @Success 200 {array} models.User
// @Router /users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Limit results to prevent performance issues
	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit := parseLimit(limitParam); parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	// Search users by username or email (case insensitive)
	var users []models.User
	searchPattern := "%" + strings.ToLower(query) + "%"
	
	err := h.db.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ?", searchPattern, searchPattern).
		Select("id, username, name, email, avatar_url, created_at").
		Limit(limit).
		Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Get all users
// @Description Get all users (admin only or for development)
// @Tags users
// @Security BearerAuth
// @Param limit query int false "Limit results (default: 20, max: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Parse pagination parameters
	limit := 20
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit := parseLimit(limitParam); parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if parsedOffset := parseLimit(offsetParam); parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	var users []models.User
	err := h.db.Select("id, username, name, email, avatar_url, created_at").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Helper function to parse limit parameter
func parseLimit(limitStr string) int {
	var limit int
	if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil {
		return 0
	}
	return limit
}