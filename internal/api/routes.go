package api

import (
	"devsync-be/internal/api/handlers"
	"devsync-be/internal/storage"
    "devsync-be/internal/api/middleware"
    "devsync-be/internal/config"
    "devsync-be/internal/websocket"

    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, hub *websocket.Hub, cfg *config.Config, gcsStorage *storage.GCSStorage) {
    // Middleware
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(middleware.CORS())

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(db, cfg)
    projectHandler := handlers.NewProjectHandler(db)
    fileHandler := handlers.NewFileHandler(db, hub)
    uploadHandler := handlers.NewUploadHandler(db, gcsStorage)
    taskHandler := handlers.NewTaskHandler(db, hub)
    chatHandler := handlers.NewChatHandler(db, hub)

    // Swagger documentation
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // WebSocket endpoint
    r.GET("/ws", hub.HandleWebSocket)

    // API routes
    api := r.Group("/api/v1")
    {
        // Auth routes
        auth := api.Group("/auth")
        {
            auth.GET("/github", authHandler.GitHubLogin)
            auth.GET("/github/callback", authHandler.GitHubCallback)
            auth.POST("/dev-login", authHandler.DevLogin) // Tambahkan ini
            auth.POST("/refresh", authHandler.RefreshToken)
        }

        // Protected routes
        protected := api.Group("/")
        protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
        {
            // User routes
            protected.GET("/me", authHandler.GetCurrentUser)

            // Project routes
            projects := protected.Group("/projects")
            {
                projects.GET("/", projectHandler.GetProjects)
                projects.POST("/", projectHandler.CreateProject)
                projects.GET("/:id", projectHandler.GetProject)
                projects.PUT("/:id", projectHandler.UpdateProject)
                projects.DELETE("/:id", projectHandler.DeleteProject)

                // File routes
                projects.GET("/:id/files", fileHandler.GetFiles)
                projects.POST("/:id/files", fileHandler.CreateFile)
                projects.GET("/:id/files/:fileId", fileHandler.GetFile)
                projects.PUT("/:id/files/:fileId", fileHandler.UpdateFile)
                projects.DELETE("/:id/files/:fileId", fileHandler.DeleteFile)

                // Task routes
                projects.GET("/:id/tasks", taskHandler.GetTasks)
                projects.POST("/:id/tasks", taskHandler.CreateTask)
                projects.PUT("/:id/tasks/:taskId", taskHandler.UpdateTask)
                projects.DELETE("/:id/tasks/:taskId", taskHandler.DeleteTask)

                // Sprint routes
                projects.GET("/:id/sprints", taskHandler.GetSprints)
                projects.POST("/:id/sprints", taskHandler.CreateSprint)

                // Chat routes
                projects.GET("/:id/messages", chatHandler.GetMessages)
                projects.POST("/:id/messages", chatHandler.SendMessage)
                projects.POST("/:id/upload", uploadHandler.UploadFile)
            }
        }
    }
}