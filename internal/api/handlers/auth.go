package handlers

import (
    "net/http"

    "devsync-be/internal/auth"
    "devsync-be/internal/config"
    "devsync-be/internal/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type AuthHandler struct {
    db  *gorm.DB
    cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
    return &AuthHandler{
        db:  db,
        cfg: cfg,
    }
}

// @Summary GitHub OAuth login
// @Description Redirect to GitHub OAuth
// @Tags auth
// @Success 302 {string} string "redirect"
// @Router /auth/github [get]
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
    config := auth.GetGitHubOAuthConfig(h.cfg.GitHubClientID, h.cfg.GitHubSecret, h.cfg.RedirectURL)
    url := config.AuthCodeURL("state")
    c.Redirect(http.StatusTemporaryRedirect, url)
}

// @Summary GitHub OAuth callback
// @Description Handle GitHub OAuth callback
// @Tags auth
// @Param code query string true "OAuth code"
// @Success 200 {object} map[string]interface{}
// @Router /auth/github/callback [get]
func (h *AuthHandler) GitHubCallback(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Code not provided"})
        return
    }

    config := auth.GetGitHubOAuthConfig(h.cfg.GitHubClientID, h.cfg.GitHubSecret, h.cfg.RedirectURL)
    githubUser, accessToken, err := auth.GetGitHubUser(c.Request.Context(), config, code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
        return
    }

    // Find or create user - gunakan git_hub_id bukan github_id
    var user models.User
    result := h.db.Where("git_hub_id = ?", githubUser.ID).First(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            // Create new user
            user = models.User{
                GitHubID:    githubUser.ID,
                Username:    githubUser.Login,
                Email:       githubUser.Email,
                Name:        githubUser.Name,
                AvatarURL:   githubUser.AvatarURL,
                AccessToken: accessToken,
            }
            if createErr := h.db.Create(&user).Error; createErr != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
                return
            }
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            return
        }
    } else {
        // Update existing user
        user.AccessToken = accessToken
        if updateErr := h.db.Save(&user).Error; updateErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
            return
        }
    }

    user.AvatarURL = githubUser.AvatarURL
    user.Name = githubUser.Name
    h.db.Save(&user)

    // Generate JWT token
    token, err := auth.GenerateToken(user.ID, user.Username, h.cfg.JWTSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Return success response with token and user info
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "token":   token,
        "user": gin.H{
            "id":         user.ID,
            "username":   user.Username,
            "email":      user.Email,
            "name":       user.Name,
            "avatar_url": user.AvatarURL,
        },
    })
}

// @Summary Refresh token
// @Description Refresh JWT token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    username, exists := c.Get("username")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
        return
    }

    token, err := auth.GenerateToken(userID.(uint), username.(string), h.cfg.JWTSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Get current user
// @Description Get current authenticated user
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} models.User
// @Router /me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    var user models.User
    if err := h.db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}