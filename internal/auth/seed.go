package auth

import (
	"context"
	"log"

	"github.com/karanbihani/file-vault/internal/db" // Adjust path
	"github.com/jackc/pgx/v5"
)

// SeedRBAC ensures that the default roles and permissions exist in the database.
func SeedRBAC(ctx context.Context, queries *db.Queries) {
	log.Println("Seeding RBAC roles and permissions...")

	// --- Define Roles and their Permissions ---
	rolesAndPermissions := map[string][]string{
		"user": {
			"files:upload",
			"files:read:self",
			"files:delete:self",
			"shares:create:self",
			"stats:read:self",
		},
		"admin": {
			"files:read:all",
			"files:delete:any",
			"stats:read:all",
		},
	}

	// --- Create Roles and Permissions, and Link them ---
	for roleName, permissions := range rolesAndPermissions {
		// Get or create the role
		role, err := queries.GetRoleByName(ctx, roleName)
		if err == pgx.ErrNoRows {
			role, err = queries.CreateRole(ctx, roleName)
			if err != nil {
				log.Fatalf("Failed to create role '%s': %v", roleName, err)
			}
		} else if err != nil {
			log.Fatalf("Failed to get role '%s': %v", roleName, err)
		}

		// For each permission for this role...
		for _, permName := range permissions {
			// Get or create the permission
			permission, err := queries.GetPermissionByName(ctx, permName)
			if err == pgx.ErrNoRows {
				permission, err = queries.CreatePermission(ctx, permName)
				if err != nil {
					log.Fatalf("Failed to create permission '%s': %v", permName, err)
				}
			} else if err != nil {
				log.Fatalf("Failed to get permission '%s': %v", permName, err)
			}

			// Link the role to the permission
			err = queries.LinkRoleToPermission(ctx, db.LinkRoleToPermissionParams{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			})
			// We can ignore duplicate key errors, as it just means the link already exists.
			if err != nil && !isDuplicateKeyError(err) {
				log.Fatalf("Failed to link role '%s' to permission '%s': %v", roleName, permName, err)
			}
		}
	}
	log.Println("RBAC seeding complete.")
}

// isDuplicateKeyError is a helper to check for unique constraint violation errors.
func isDuplicateKeyError(err error) bool {
	// This is specific to pgx errors.
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505"
	}
}