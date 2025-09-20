package auth

// Defines the set of granular permissions available in the system.
const (
    // File Permissions
    PermissionFilesUpload     = "files:upload"
    PermissionFilesDownload   = "files:download"
    PermissionFilesDelete     = "files:delete"
    PermissionFilesReadShared = "files:read:shared" // For 'shared-with-me'

    // Sharing Permissions (These are the ones we will use for sharing actions)
    PermissionSharesCreatePublic = "shares:create:public"
    PermissionSharesCreateUser   = "shares:create:user"
    PermissionSharesRevokePublic = "shares:revoke:public"
    PermissionSharesRevokeUser   = "shares:revoke:user"

    // Statistics Permissions
    PermissionStatsReadSelf = "stats:read:self"

    // Admin-level permissions
    PermissionAdminManageRoles = "admin:manage_roles"
    PermissionAdminViewAllFiles = "admin:view_all_files" // <-- ADD THIS
    PermissionAdminViewAllStats = "admin:view_all_stats" // <-- ADD THIS
    PermissionAdminDownloadAnyFile = "admin:download_any_file" // <-- ADD THIS
)