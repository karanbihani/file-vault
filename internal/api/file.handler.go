package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/karanbihani/file-vault/internal/core/files" // Adjust path
	"github.com/gin-gonic/gin"
)

// ... (FilesHandler struct and NewFilesHandler are the same)
type FilesHandler struct {
	fileService *files.Service
}

func NewFilesHandler(service *files.Service) *FilesHandler {
	return &FilesHandler{
		fileService: service,
	}
}

// Upload now gets the ownerID from the context.
func (h *FilesHandler) Upload(c *gin.Context) {
	// Get the user ID from the context that the AuthMiddleware set.
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required in the 'file' form field"})
		return
	}
	defer file.Close()

	description := c.PostForm("description")
	tags := c.PostFormArray("tags")
	if len(tags) == 1 && strings.Contains(tags[0], ",") {
		tags = strings.Split(tags[0], ",")
	}

	uploadParams := files.UploadFileParams{
		File:        file,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		OwnerID:     userID.(int64), // We assert the type to int64
		Description: description,
		Tags:        tags,
	}

	userFile, err := h.fileService.UploadFile(c.Request.Context(), uploadParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userFile)
}

// List now gets the ownerID from the context.
func (h *FilesHandler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	files, err := h.fileService.ListFiles(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

// Download now gets the ownerID from the context.
func (h *FilesHandler) Download(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	downloadData, err := h.fileService.DownloadFile(c.Request.Context(), fileID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer downloadData.Data.Close()

	c.Header("Content-Disposition", "attachment; filename="+downloadData.Filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", downloadData.Size))
	c.DataFromReader(http.StatusOK, downloadData.Size, "application/octet-stream", downloadData.Data, nil)
}

// Delete now gets the ownerID from the context.
func (h *FilesHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	err = h.fileService.DeleteFile(c.Request.Context(), fileID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})
}

// ListSharedWithMe is the handler for the GET /files/shared-with-me endpoint.
func (h *FilesHandler) ListSharedWithMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	files, err := h.fileService.ListFilesSharedWithMe(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}