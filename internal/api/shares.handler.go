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