package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Factory function to create a handler with access to the database pool.
func HealthCheckHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Starting health check")
		if dbpool == nil {
			log.Println("Database pool is nil")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "DOWN", "error": "database pool is not initialized"})
			return
		}

		// Ping the database to ensure the connection is alive.
		err := dbpool.Ping(context.Background())
		if err != nil {
			log.Printf("Database ping failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "DOWN", "error": "database connection error"})
			return
		}

		log.Println("Health check passed")
		c.JSON(http.StatusOK, gin.H{"status": "UP", "database": "healthy"})
	}
}