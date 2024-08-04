package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/habibaabulata/file_storage/handlers"
	"github.com/habibaabulata/file_storage/storage"
)

func main() {
	// Initialize the database
	storage.InitDB("C:/Users/habib/OneDrive/Desktop/workspace/fstorage/server/path/to/leveldb")
	defer storage.CloseDB()

	// Initialize the Gin router
	router := gin.Default()

	// Public route for login
	router.POST("/login", handlers.Login)

	// Group for authorized routes
	authorized := router.Group("/")
	authorized.Use(handlers.AuthMiddleware())
	{
		authorized.POST("/upload", handlers.UploadFile)
		authorized.GET("/download/:fileID", handlers.DownloadFile)
	}

	// Use TLS for HTTPS
	if err := router.RunTLS(":8080", "C:/Users/habib/OneDrive/Desktop/workspace/fstorage/cert.pem", "C:/Users/habib/OneDrive/Desktop/workspace/fstorage/key.pem"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
