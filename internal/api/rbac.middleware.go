package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/karanbihani/file-vault/internal/db"
)

// PermissionMiddleware is a factory that creates a Gin middleware to enforce a required permission.
// It queries the database on each request to get the user's current permissions.
func PermissionMiddleware(queries *db.Queries, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		permissions, err := queries.GetUserPermissions(c.Request.Context(), userID.(int64))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve user permissions"})
			return
		}

		for _, p := range permissions {
			if p == requiredPermission {
				c.Next()
				return
			}
		}

		errorMsg := gin.H{"error": "access denied: you do not have the required permission (" + requiredPermission + ")"}
		c.AbortWithStatusJSON(http.StatusForbidden, errorMsg)
	}
}