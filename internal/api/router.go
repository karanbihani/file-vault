package api

import (
	"github.com/karanbihani/file-vault/internal/core/files"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRouter now accepts the FileService as a dependency.
func SetupRouter(dbpool *pgxpool.Pool, fileService *files.Service) *gin.Engine {
	router := gin.Default()

	// --- Handler Initialization ---
	// We create an instance of our FilesHandler, injecting the fileService.
	fileHandler := NewFilesHandler(fileService)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", HealthCheckHandler(dbpool))

		// --- File Routes ---
		// We register the new upload route and connect it to our handler method.
		v1.POST("/files", fileHandler.Upload)

		// We will add GET, DELETE etc. here in the next task.
	}

	return router
}