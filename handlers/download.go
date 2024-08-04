package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibaabulata/file_storage/storage"
)

func DownloadFile(c *gin.Context) {
	fileID := c.Param("fileID")
	log.Printf("Received request for fileID: %s", fileID)

	// Retrieve metadata
	metadataBytes, err := storage.DB.Get([]byte(fileID), nil)
	if err != nil {
		log.Printf("Error retrieving metadata from database for fileID %s: %v", fileID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	log.Printf("Successfully retrieved metadata for file ID: %s", fileID)

	var metadata map[string]interface{}
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		log.Printf("Error unmarshalling metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file metadata"})
		return
	}

	chunkCount, ok := metadata["chunkCount"].(float64)
	if !ok {
		log.Println("Invalid chunkCount in metadata")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid file metadata"})
		return
	}

	var fileContent []byte
	for i := 0; i < int(chunkCount); i++ {
		chunkPath := filepath.Join("uploads", fileID+"_"+strconv.Itoa(i))
		log.Printf("Reading chunk from path: %s", chunkPath)
		chunkContent, err := ioutil.ReadFile(chunkPath)
		if err != nil {
			log.Printf("Error reading file chunk %d: %v", i, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file chunk"})
			return
		}
		fileContent = append(fileContent, chunkContent...)
	}

	log.Printf("Sending file with size: %d bytes", len(fileContent))
	// Send the file content as the response
	c.Header("Content-Disposition", "attachment; filename="+metadata["filename"].(string))
	c.Data(http.StatusOK, "application/octet-stream", fileContent)
}
