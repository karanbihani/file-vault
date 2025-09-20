package rbac

import (
	"context"
	"fmt"

	"github.com/karanbihani/file-vault/internal/db"
)

// Service handles the business logic for RBAC.
type Service struct {
	queries *db.Queries
}

// NewService creates a new RBAC service.
func NewService(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

// ListRoles retrieves all roles from the database.
func (s *Service) ListRoles(ctx context.Context) ([]db.Role, error) {
	return s.queries.ListRoles(ctx)
}

// ListPermissions retrieves all permissions from the database.
func (s *Service) ListPermissions(ctx context.Context) ([]db.Permission, error) {
	return s.queries.ListPermissions(ctx)
}

// GetPermissionsForRole retrieves all permissions for a specific role ID.
func (s *Service) GetPermissionsForRole(ctx context.Context, roleID int32) ([]db.Permission, error) {
	return s.queries.GetPermissionsForRole(ctx, roleID)
}

// AddPermissionToRole assigns a permission to a role.
func (s *Service) AddPermissionToRole(ctx context.Context, roleID, permissionID int32) error {
	err := s.queries.AddPermissionToRole(ctx, db.AddPermissionToRoleParams{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
	if err != nil {
		return fmt.Errorf("could not add permission to role: %w", err)
	}
	return nil
}

// RemovePermissionFromRole removes a permission from a role.
func (s *Service) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int32) error {
	err := s.queries.RemovePermissionFromRole(ctx, db.RemovePermissionFromRoleParams{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
	if err != nil {
		return fmt.Errorf("could not remove permission from role: %w", err)
	}
	return nil
}