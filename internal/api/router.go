package api

import (
	"time"

	"github.com/karanbihani/file-vault/internal/auth"    // Adjust path
	"github.com/karanbihani/file-vault/internal/core/files" // Adjust path
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(dbpool *pgxpool.Pool, fileService *files.Service, authService *auth.Service) *gin.Engine {
	router := gin.Default()

	fileHandler := NewFilesHandler(fileService)
	authHandler := NewAuthHandler(authService)

	router.Use(RateLimiter(2, time.Second))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", HealthCheckHandler(dbpool))

		// --- Public Auth Routes ---
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)

		// --- Protected File Routes ---
		// We create a new group for routes that require authentication.
		protected := v1.Group("/")
		// We apply our AuthMiddleware to this entire group.
		protected.Use(AuthMiddleware())
		{
			protected.POST("/files", fileHandler.Upload)
			protected.GET("/files", fileHandler.List)
			protected.GET("/files/:id/download", fileHandler.Download)
			protected.DELETE("/files/:id", fileHandler.Delete)
		}
	}
	return router
}