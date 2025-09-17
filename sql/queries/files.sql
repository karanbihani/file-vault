-- name: GetPhysicalFileByHash :one
-- Retrieves a single physical_file record by its SHA-26 hash to check for duplicates.
SELECT * FROM physical_files
WHERE sha256_hash = $1 LIMIT 1;

-- name: CreatePhysicalFile :one
-- Inserts a new physical_file record into the database when a new unique file is uploaded.
-- It returns the newly created record.
INSERT INTO physical_files (
  sha256_hash,
  size_bytes,
  storage_path
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: IncrementPhysicalFileRefCount :one
-- Increments the reference_count for a physical_file when a duplicate is uploaded.
-- It returns the updated record.
UPDATE physical_files
SET reference_count = reference_count + 1
WHERE id = $1
RETURNING *;

-- name: CreateUserFile :one
-- Inserts a new user_file record, linking a user to a physical file.
-- UPDATED: Now includes description and tags.
INSERT INTO user_files (
  owner_id,
  physical_file_id,
  filename,
  mime_type,
  description,
  tags
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;