package api

import (
	"net/http"
	"strings"

	"github.com/karanbihani/file-vault/internal/core/files" 
	"github.com/gin-gonic/gin"
)

type FilesHandler struct {
	fileService *files.Service
}

func NewFilesHandler(service *files.Service) *FilesHandler {
	return &FilesHandler{
		fileService: service,
	}
}

func (h *FilesHandler) Upload(c *gin.Context) {
	// For now, we will hardcode the ownerID. In Day 3, we will get this
	// from the JWT token after the user logs in.
	const ownerID = 1

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required in the 'file' form field"})
		return
	}
	defer file.Close()

	// --- ADDED: Read optional description and tags from the form data ---
	// c.PostForm reads a single value for a key.
	description := c.PostForm("description")

	// c.PostFormArray can read multiple values for the same key (e.g., tags=tag1&tags=tag2).
	// We also handle a single comma-separated string for tags for flexibility.
	tags := c.PostFormArray("tags")
	if len(tags) == 1 && strings.Contains(tags[0], ",") {
		tags = strings.Split(tags[0], ",")
	}
	// --- END ADDITION ---

	// Prepare the parameters for our service method.
	uploadParams := files.UploadFileParams{
		File:        file,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		OwnerID:     ownerID,
		Description: description, // Pass the new data
		Tags:        tags,        // Pass the new data
	}

	// Call the core business logic in the service.
	userFile, err := h.fileService.UploadFile(c.Request.Context(), uploadParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userFile)
}
