package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/karanbihani/file-vault/internal/core/admin"
)

type AdminHandler struct {
	adminService *admin.Service
}

func NewAdminHandler(service *admin.Service) *AdminHandler {
	return &AdminHandler{adminService: service}
}

func (h *AdminHandler) ListAllFiles(c *gin.Context) {
	files, err := h.adminService.ListAllFiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

func (h *AdminHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.adminService.GetSystemStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) ListAuditLogs(c *gin.Context) {
	logs, err := h.adminService.ListAuditLogs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
