package handlers

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/habibaabulata/file_storage/storage" // Adjust import as needed
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadFile(t *testing.T) {
	// Initialize the database
	storage.InitDB("C:/Users/habib/OneDrive/Desktop/workspace/fstorage/server/path/to/leveldb")
	defer storage.CloseDB()

	// Create the uploads directory if it doesn't exist
	err := os.MkdirAll(
		"uploads",
		os.ModePerm,
	)
	require.NoError(t, err)

	// Initialize the Gin router
	router := gin.Default()
	router.POST(
		"/upload",
		UploadFile,
	)

	// Create a buffer to hold the multipart form data
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Create a form file field
	fileWriter, err := writer.CreateFormFile(
		"file",
		"testfile.txt",
	)
	require.NoError(t, err)

	// Write some content to the form file field
	_, err = fileWriter.Write([]byte("This is a test file"))
	require.NoError(t, err)

	// Close the writer to set the ending boundary
	err = writer.Close()
	require.NoError(t, err)

	// Create a new HTTP request with the multipart form data
	req := httptest.NewRequest(
		http.MethodPost,
		"/upload",
		&buffer,
	)
	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	// Create a response recorder to capture the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the file was saved correctly
	savedFiles, err := ioutil.ReadDir("uploads")
	require.NoError(t, err)
	require.NotEmpty(t, savedFiles)

	// Clean up
	for _, f := range savedFiles {
		os.Remove(filepath.Join(
			"uploads",
			f.Name(),
		))
	}
}
