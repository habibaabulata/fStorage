package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/habibaabulata/file_storage/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Set up LevelDB
	storage.InitDB("C:/Users/habib/OneDrive/Desktop/workspace/fstorage/server/path/to/leveldb")
	defer storage.CloseDB()

	// Run tests
	code := m.Run()

	// Clean up
	os.Exit(code)
}

func TestDownloadFile(t *testing.T) {
	// Create the uploads directory if it doesn't exist
	err := os.MkdirAll("uploads", os.ModePerm)
	require.NoError(t, err)

	// Create a temporary test file
	testFilePath := "uploads/testfile.txt"
	err = ioutil.WriteFile(testFilePath, []byte("This is a test file"), 0644)
	require.NoError(t, err)

	// Store the test file path in the database
	fileID := "testfileid"
	err = storage.DB.Put([]byte(fileID), []byte(testFilePath), nil)
	require.NoError(t, err)

	// Initialize the Gin router
	router := gin.Default()
	router.GET("/download/:fileID", DownloadFile)

	// Create a new HTTP request to download the file
	req := httptest.NewRequest(http.MethodGet, "/download/"+fileID, nil)

	// Create a response recorder to capture the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	responseBody, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	assert.Equal(t, "This is a test file", string(responseBody))

	// Clean up the test file
	os.Remove(testFilePath)
}
