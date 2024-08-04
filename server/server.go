package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/habibaabulata/file_storage/handlers"
	"github.com/habibaabulata/file_storage/storage"
)

func main() {
	// Initialize the database with the path to your LevelDB
	storage.InitDB("C:/Users/habib/OneDrive/Desktop/workspace/fstorage/server/path/to/leveldb")
	defer storage.CloseDB()

	// Initialize the Gin router
	router := gin.Default()

	// Define routes
	router.POST("/login", handlers.Login)
	authorized := router.Group("/")
	authorized.Use(handlers.AuthMiddleware())
	{
		authorized.POST("/upload", handlers.UploadFile)
		authorized.GET("/download/:fileID", handlers.DownloadFile)
	}

	// Use TLS for HTTPS
	if err := router.RunTLS(":8080", "cert.pem", "key.pem"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
