package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/karanbihani/file-vault/internal/core/shares" // Adjust path
	"github.com/gin-gonic/gin"
)

type SharesHandler struct {
	sharesService *shares.Service
}

func NewSharesHandler(service *shares.Service) *SharesHandler {
	return &SharesHandler{
		sharesService: service,
	}
}

// CreatePublicLink is the PROTECTED handler for POST /files/:id/share
func (h *SharesHandler) CreatePublicLink(c *gin.Context) {
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

	share, err := h.sharesService.CreatePublicLink(c.Request.Context(), fileID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the full public URL to the user.
	// NOTE: In production, you would use a frontend URL, not the API host.
	fullURL := fmt.Sprintf("http://%s/api/v1/share/%s", c.Request.Host, share.ShareToken)
	c.JSON(http.StatusOK, gin.H{"share_url": fullURL})
}

// PublicDownload is the PUBLIC handler for GET /share/:token
func (h *SharesHandler) PublicDownload(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "share token is required"})
		return
	}

	downloadData, err := h.sharesService.ProcessPublicDownload(c.Request.Context(), token)
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

func (h *SharesHandler) ShareWithUser(c *gin.Context) {
	// Define a struct to bind the incoming JSON request body.
	var requestBody struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: 'email' field is required"})
		return
	}

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

	err = h.sharesService.ShareFileWithUser(c.Request.Context(), fileID, userID.(int64), requestBody.Email)
	if err != nil {
		// We can check for specific error messages to return better status codes in the future.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("file successfully shared with %s", requestBody.Email)})
}

func (h *SharesHandler) RevokePublicLinks(c *gin.Context) {
	userID, _ := c.Get("userID")
	fileID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.sharesService.RevokePublicLinks(c.Request.Context(), fileID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all public share links for the file have been revoked"})
}

// UnshareWithUser is the PROTECTED handler for DELETE /files/:id/share-to-user
func (h *SharesHandler) UnshareWithUser(c *gin.Context) {
	var requestBody struct {
		RecipientID int64 `json:"recipient_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: 'recipient_id' is required"})
		return
	}

	userID, _ := c.Get("userID")
	fileID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.sharesService.UnshareFileWithUser(c.Request.Context(), fileID, userID.(int64), requestBody.RecipientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "share access has been revoked for the specified user"})
}