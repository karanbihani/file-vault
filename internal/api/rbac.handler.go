package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/karanbihani/file-vault/internal/core/rbac"
)

type RBACHandler struct {
	rbacService *rbac.Service
}

func NewRBACHandler(service *rbac.Service) *RBACHandler {
	return &RBACHandler{
		rbacService: service,
	}
}

func (h *RBACHandler) ListRoles(c *gin.Context) {
	roles, err := h.rbacService.ListRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *RBACHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.rbacService.ListPermissions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

func (h *RBACHandler) GetPermissionsForRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	permissions, err := h.rbacService.GetPermissionsForRole(c.Request.Context(), int32(roleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

func (h *RBACHandler) AddPermissionToRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}
	permissionID, err := strconv.ParseInt(c.Param("permissionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}

	err = h.rbacService.AddPermissionToRole(c.Request.Context(), int32(roleID), int32(permissionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission added to role successfully"})
}

func (h *RBACHandler) RemovePermissionFromRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}
	permissionID, err := strconv.ParseInt(c.Param("permissionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}

	err = h.rbacService.RemovePermissionFromRole(c.Request.Context(), int32(roleID), int32(permissionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission removed from role successfully"})
}