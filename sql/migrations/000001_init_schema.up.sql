-- This schema defines the complete database structure for the File Vault application.
-- It is designed for PostgreSQL and includes tables for users, RBAC,
-- file storage with deduplication, and file sharing.

-- Users table stores authentication and quota information.
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    storage_quota_bytes BIGINT NOT NULL DEFAULT 1073741824, -- Default 1 GB
    storage_used_bytes BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Roles table for Role-Based Access Control (RBAC).
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL -- e.g., 'user', 'admin'
);

-- Permissions table for granular RBAC.
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL -- e.g., 'files:upload', 'admin:view_all'
);

-- Junction table to link users to roles (many-to-many).
CREATE TABLE user_roles (
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role_id INT REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Junction table to link roles to permissions (many-to-many).
CREATE TABLE role_permissions (
    role_id INT REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INT REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Core table for deduplication. Stores the actual physical file info.
CREATE TABLE physical_files (
    id BIGSERIAL PRIMARY KEY,
    sha256_hash VARCHAR(64) UNIQUE NOT NULL,
    size_bytes BIGINT NOT NULL,
    storage_path TEXT NOT NULL,
    reference_count INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index on the hash for fast lookups during upload to check for duplicates.
CREATE INDEX idx_physical_files_hash ON physical_files(sha256_hash);

-- User-facing file metadata. Links a user to a physical file.
CREATE TABLE user_files (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    physical_file_id BIGINT NOT NULL REFERENCES physical_files(id) ON DELETE RESTRICT,
    filename VARCHAR(255) NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    upload_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table to manage public sharing links.
CREATE TABLE shares (
    id BIGSERIAL PRIMARY KEY,
    user_file_id BIGINT NOT NULL REFERENCES user_files(id) ON DELETE CASCADE,
    share_token VARCHAR(32) UNIQUE NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    download_count BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- *** NEW: Table for sharing files with specific users (Bonus Feature) ***
CREATE TABLE file_shares_to_users (
    user_file_id BIGINT NOT NULL REFERENCES user_files(id) ON DELETE CASCADE,
    shared_with_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_file_id, shared_with_user_id)
);

-- *** NEW: Table for audit logging (Bonus Feature) ***
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL, -- Keep log even if user is deleted
    action VARCHAR(255) NOT NULL,
    details JSONB,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);