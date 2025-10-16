package main

import (
	"log"
	"os"

	"devsync-be/internal/api"
	"devsync-be/internal/config"
	"devsync-be/internal/database"
	"devsync-be/internal/storage"
	"devsync-be/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize GCS storage
	gcsStorage, err := storage.NewGCSStorage(cfg.GCPProjectID, cfg.GCPBucketName, cfg.GCPCredentialsPath)
	if err != nil {
		log.Fatal("Failed to initialize GCS storage:", err)
	}
	defer gcsStorage.Close()

	hub := websocket.NewHub(cfg)
	go hub.Run()

	r := gin.Default()

	api.SetupRoutes(r, db, hub, cfg, gcsStorage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}