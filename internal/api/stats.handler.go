package api

import (
	"net/http"

	"github.com/karanbihani/file-vault/internal/core/stats" // Adjust to your module path
	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *stats.Service
}

func NewStatsHandler(service *stats.Service) *StatsHandler {
	return &StatsHandler{
		statsService: service,
	}
}

// GetUserDashboardStats is the PROTECTED handler for GET /stats
func (h *StatsHandler) GetUserDashboardStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	userStats, err := h.statsService.GetUserDashboardStats(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userStats)
}