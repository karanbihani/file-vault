package api

import (
	"time"

	"github.com/karanbihani/file-vault/internal/auth"    // Adjust path
	"github.com/gin-contrib/cors"
	"github.com/karanbihani/file-vault/internal/core/files" // Adjust path
	"github.com/karanbihani/file-vault/internal/core/shares" // Add this import
	"github.com/karanbihani/file-vault/internal/core/stats"  // Add this import
	"github.com/karanbihani/file-vault/internal/core/rbac" 
	"github.com/karanbihani/file-vault/internal/db" // <-- Add this import for db.Queries
	"github.com/karanbihani/file-vault/internal/core/admin" // <-- Add this import for admin service
	"github.com/karanbihani/file-vault/internal/core/search"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(queries *db.Queries, dbpool *pgxpool.Pool, fileService *files.Service, authService *auth.Service, sharesService *shares.Service,
	statsService *stats.Service, rbacService *rbac.Service, adminService *admin.Service, searchService *search.Service) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	fileHandler := NewFilesHandler(fileService)
	authHandler := NewAuthHandler(authService)
	sharesHandler := NewSharesHandler(sharesService)
	statsHandler := NewStatsHandler(statsService) // Create the new handler
	rbacHandler := NewRBACHandler(rbacService) // <-- Initialize the new RBAC handler
	adminHandler := NewAdminHandler(adminService) // <-- Initialize the new Admin handler
	searchHandler := NewSearchHandler(searchService) // <-- Initialize the new handler

	router.Use(RateLimiter(2, time.Second))

	v1 := router.Group("/api/v1")
	{
		// --- Public Routes ---
		v1.GET("/health", HealthCheckHandler(dbpool))
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.GET("/share/:token", sharesHandler.PublicDownload)

		// --- Protected User Routes ---
		// All routes in this group require authentication first.
		// Then, each route has a specific permission check.
		protected := v1.Group("/")
		protected.Use(AuthMiddleware())
		{
			// File Management Routes
			protected.POST("/files", PermissionMiddleware(queries, auth.PermissionFilesUpload), fileHandler.Upload)
			protected.GET("/files", fileHandler.List) // Listing own files doesn't need a specific perm
			protected.GET("/files/:id/download", PermissionMiddleware(queries, auth.PermissionFilesDownload), fileHandler.Download)
			protected.DELETE("/files/:id", PermissionMiddleware(queries, auth.PermissionFilesDelete), fileHandler.Delete)
			protected.GET("/files/shared-with-me", PermissionMiddleware(queries, auth.PermissionFilesReadShared), fileHandler.ListSharedWithMe) // Assuming List handler can be adapted

			// Sharing Management Routes
			protected.POST("/files/:id/share", PermissionMiddleware(queries, auth.PermissionSharesCreatePublic), sharesHandler.CreatePublicLink)
			protected.POST("/files/:id/share-to-user", PermissionMiddleware(queries, auth.PermissionSharesCreateUser), sharesHandler.ShareWithUser)
			protected.DELETE("/files/:id/share", PermissionMiddleware(queries, auth.PermissionSharesRevokePublic), sharesHandler.RevokePublicLinks)
			protected.DELETE("/files/:id/share-to-user", PermissionMiddleware(queries, auth.PermissionSharesRevokeUser), sharesHandler.UnshareWithUser)
			protected.GET("/files/:id/shares", sharesHandler.GetSharesForFile) 

			// Stats Route
			protected.GET("/stats", PermissionMiddleware(queries, auth.PermissionStatsReadSelf), statsHandler.GetUserDashboardStats)

			// Search Route
			protected.GET("/search", searchHandler.Search)

			// Tag Management Route
			protected.POST("/files/:id/tags", fileHandler.AddTag)
			protected.DELETE("/files/:id/tags", fileHandler.RemoveTag)
		}


		// --- Protected Admin & RBAC Management Routes ---
		admin := v1.Group("/admin")
		admin.Use(AuthMiddleware())
		{
			// RBAC Management APIs
			admin.GET("/roles", PermissionMiddleware(queries, auth.PermissionAdminManageRoles), rbacHandler.ListRoles)
			
			admin.GET("/permissions", PermissionMiddleware(queries, auth.PermissionAdminManageRoles), rbacHandler.ListPermissions)
			admin.GET("/roles/:roleId/permissions", PermissionMiddleware(queries, auth.PermissionAdminManageRoles), rbacHandler.GetPermissionsForRole)
			admin.POST("/roles/:roleId/permissions/:permissionId", PermissionMiddleware(queries, auth.PermissionAdminManageRoles), rbacHandler.AddPermissionToRole)
			admin.DELETE("/roles/:roleId/permissions/:permissionId", PermissionMiddleware(queries, auth.PermissionAdminManageRoles), rbacHandler.RemovePermissionFromRole)
			
			admin.GET("/files", PermissionMiddleware(queries, auth.PermissionAdminViewAllFiles), adminHandler.ListAllFiles)
			admin.GET("/stats", PermissionMiddleware(queries, auth.PermissionAdminViewAllStats), adminHandler.GetSystemStats)
			admin.GET("/logs", PermissionMiddleware(queries, auth.PermissionAdminViewAuditLogs), adminHandler.ListAuditLogs)
		}
	}
	return router
}