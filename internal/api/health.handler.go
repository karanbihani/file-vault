package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Factory function to create a handler with access to the database pool.
func HealthCheckHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ping the database to ensure the connection is alive.
		err := dbpool.Ping(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "DOWN", "error": "database connection error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "UP", "database": "healthy"})
	}
}