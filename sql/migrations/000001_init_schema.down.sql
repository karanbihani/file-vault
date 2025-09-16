-- This file is for rolling back the initial schema migration.
-- It drops all tables and indexes created in the corresponding .up.sql file.
-- The order of operations is the reverse of the creation order to respect foreign key constraints.

-- Drop tables that have foreign keys first.
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS file_shares_to_users;
DROP TABLE IF EXISTS shares;
DROP TABLE IF EXISTS user_files;

-- Drop junction tables for RBAC.
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;

-- Drop the index before the table it belongs to.
DROP INDEX IF EXISTS idx_physical_files_hash;

-- Drop the main entity tables.
DROP TABLE IF EXISTS physical_files;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;

-- Note: In a production environment, consider archiving data before dropping tables.