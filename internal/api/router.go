package api

import (
	"time"

	"github.com/karanbihani/file-vault/internal/auth"    // Adjust path
	"github.com/karanbihani/file-vault/internal/core/files" // Adjust path
	"github.com/karanbihani/file-vault/internal/core/shares" // Add this import
	"github.com/karanbihani/file-vault/internal/core/stats"  // Add this import
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(dbpool *pgxpool.Pool, fileService *files.Service, authService *auth.Service, sharesService *shares.Service, statsService *stats.Service) *gin.Engine {
	router := gin.Default()

	fileHandler := NewFilesHandler(fileService)
	authHandler := NewAuthHandler(authService)
	sharesHandler := NewSharesHandler(sharesService)
	statsHandler := NewStatsHandler(statsService) // Create the new handler

	router.Use(RateLimiter(2, time.Second))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", HealthCheckHandler(dbpool))

		// --- Public Auth Routes ---
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)

		// --- Public Share Routes ---
		v1.GET("/share/:token", sharesHandler.PublicDownload)

		// --- Protected File Routes ---
		// We create a new group for routes that require authentication.
		protected := v1.Group("/")
		protected.Use(AuthMiddleware())
		{
			protected.POST("/files", fileHandler.Upload)
			protected.GET("/files", fileHandler.List)
			protected.GET("/files/:id/download", fileHandler.Download)
			protected.DELETE("/files/:id", fileHandler.Delete)
			protected.GET("/files/shared-with-me", fileHandler.ListSharedWithMe) 

			protected.POST("/files/:id/share", sharesHandler.CreatePublicLink)
			protected.POST("/files/:id/share-to-user", sharesHandler.ShareWithUser)
			protected.DELETE("/files/:id/share", sharesHandler.RevokePublicLinks)        
			protected.DELETE("/files/:id/share-to-user", sharesHandler.UnshareWithUser)   

			protected.GET("/stats", statsHandler.GetUserDashboardStats)
		}

		// --- ADMIN PROTECTED ROUTES ---
		adminRoutes := v1.Group("/admin")
		adminRoutes.Use(AuthMiddleware()) // Admins must be logged in...
		adminRoutes.Use(AuthorizationMiddleware(authService, "admin:view_all_files")) // ...and must have this specific permission.
		{
			// We can now define admin-only endpoints here.
			// adminRoutes.GET("/files", adminFileHandler.ListAllFiles)
		}
	}
	return router
}