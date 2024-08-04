package handlers

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/habibaabulata/file_storage/storage"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Error retrieving the file from form-data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	fileID := uuid.New().String()
	chunkSize := 4 * 1024 * 1024 // 1 MB
	fileContent, err := file.Open()
	if err != nil {
		log.Println("Error opening the file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileContent.Close()

	var wg sync.WaitGroup
	buffer := bufio.NewReader(fileContent)
	chunkIndex := 0

	// Helper function to handle chunk writing
	writeChunk := func(chunkIndex int) {
		defer wg.Done()

		chunkPath := filepath.Join("uploads", fileID+"_"+strconv.Itoa(chunkIndex))
		chunkFile, err := os.Create(chunkPath)
		if err != nil {
			log.Println("Error creating chunk file:", err)
			return
		}
		defer chunkFile.Close()

		_, err = io.CopyN(chunkFile, buffer, int64(chunkSize))
		if err != nil && err != io.EOF {
			log.Println("Error writing file chunk:", err)
			return
		}
	}

	// Start chunk writing concurrently
	for {
		wg.Add(1)
		go writeChunk(chunkIndex)
		chunkIndex++
		if buffer.Buffered() == 0 {
			break
		}
	}
	wg.Wait()

	// Save metadata as before
	metadata := map[string]interface{}{
		"filename":   file.Filename,
		"chunkSize":  chunkSize,
		"chunkCount": chunkIndex,
		"uploadedAt": time.Now().Unix(),
	}
	metadataBytes, _ := json.Marshal(metadata)
	err = storage.DB.Put([]byte(fileID), metadataBytes, nil)
	if err != nil {
		log.Println("Error saving file metadata in database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fileID": fileID})
}
