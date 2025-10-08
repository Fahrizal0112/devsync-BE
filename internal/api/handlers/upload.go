package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"devsync-be/internal/models"
	"devsync-be/internal/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UploadHandler struct {
	db      *gorm.DB
	storage *storage.GCSStorage
}

func NewUploadHandler(db *gorm.DB, storage *storage.GCSStorage) *UploadHandler {
	return &UploadHandler{
		db:      db,
		storage: storage,
	}
}

// @Summary Upload file
// @Description Upload file to project
// @Tags files
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param file formData file true "File to upload"
// @Success 201 {object} models.File
// @Router /projects/{id}/upload [post]
func (h *UploadHandler) UploadFile(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large (max 10MB)"})
		return
	}

	// Determine file type
	fileType := getFileType(file.Filename)
	folder := fmt.Sprintf("projects/%d/%s", projectID, fileType)

	// Upload to GCS
	ctx := context.Background()
	fileURL, err := h.storage.UploadFile(ctx, file, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Save file info to database
	fileModel := models.File{
		Name:       file.Filename,
		Path:       filepath.Join(folder, file.Filename),
		FileURL:    fileURL,
		FileType:   fileType,
		FileSize:   file.Size,
		MimeType:   file.Header.Get("Content-Type"),
		ProjectID:  uint(projectID),
		UploadedBy: userID.(uint),
	}

	if err := h.db.Create(&fileModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file info"})
		return
	}

	// Load relationships
	h.db.Preload("Uploader").First(&fileModel, fileModel.ID)

	c.JSON(http.StatusCreated, fileModel)
}

func getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg"}
	documentExts := []string{".pdf", ".doc", ".docx", ".txt", ".md"}
	codeExts := []string{".go", ".js", ".py", ".java", ".cpp", ".c", ".html", ".css"}
	
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return "images"
		}
	}
	
	for _, docExt := range documentExts {
		if ext == docExt {
			return "documents"
		}
	}
	
	for _, codeExt := range codeExts {
		if ext == codeExt {
			return "code"
		}
	}
	
	return "others"
}