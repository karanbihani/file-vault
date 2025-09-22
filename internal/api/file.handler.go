package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/karanbihani/file-vault/internal/core/files" // Adjust path
	"github.com/karanbihani/file-vault/internal/db"
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

func (h *FilesHandler) Upload(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	// --- START OF FIX ---
	// Use MultipartForm to handle multiple files.
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form: " + err.Error()})
		return
	}
	// The "files" key can now contain multiple file parts.
	formFiles := form.File["files"]

	if len(formFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one file is required in the 'files' form field"})
		return
	}

	var uploadedFiles []db.UserFile

	// Loop through each file from the form.
	for _, header := range formFiles {
		file, err := header.Open()
		if err != nil {
			// Skip this file and continue with the next ones.
			log.Printf("ERROR: could not open file %s: %v", header.Filename, err)
			continue
		}
		defer file.Close()

		description := c.PostForm("description") // Description will be the same for all files in the batch
		tags := c.PostFormArray("tags")
		if len(tags) == 1 && strings.Contains(tags[0], ",") {
			tags = strings.Split(tags[0], ",")
		}

		uploadParams := files.UploadFileParams{
			File:        file,
			Filename:    header.Filename,
			ContentType: header.Header.Get("Content-Type"),
			OwnerID:     userID.(int64),
			Description: description,
			Tags:        tags,
		}

		userFile, err := h.fileService.UploadFile(c.Request.Context(), uploadParams)
		if err != nil {
			// If one file fails, we can decide to stop or continue.
			// Here, we'll log the error and continue with the other files.
			log.Printf("ERROR: failed to upload file %s: %v", header.Filename, err)
			continue
		}
		uploadedFiles = append(uploadedFiles, *userFile)
	}
	// --- END OF FIX ---

	c.JSON(http.StatusOK, uploadedFiles)
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
// AddTag handles POST /files/:id/tags
func (h *FilesHandler) AddTag(c *gin.Context) {
   // Get authenticated user ID
   userID, exists := c.Get("userID")
   if !exists {
	   c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
	   return
   }
   // Parse file ID from URL
   fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
   if err != nil {
	   c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
	   return
   }
   // Bind tag from request body
   var body struct { Tag string `json:"tag"` }
   if err := c.ShouldBindJSON(&body); err != nil {
	   c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	   return
   }
   // Call service
   if err := h.fileService.AddTag(c.Request.Context(), fileID, userID.(int64), body.Tag); err != nil {
	   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	   return
   }
   c.JSON(http.StatusOK, gin.H{"message": "tag added successfully"})
}

// RemoveTag handles DELETE /files/:id/tags
func (h *FilesHandler) RemoveTag(c *gin.Context) {
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
   var body struct { Tag string `json:"tag"` }
   if err := c.ShouldBindJSON(&body); err != nil {
	   c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	   return
   }
   if err := h.fileService.RemoveTag(c.Request.Context(), fileID, userID.(int64), body.Tag); err != nil {
	   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	   return
   }
   c.JSON(http.StatusOK, gin.H{"message": "tag removed successfully"})
}