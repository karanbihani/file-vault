-- ... (all queries up to GetFileOwnerAndPhysicalFile are the same) ...
-- name: GetPhysicalFileByHash :one
SELECT * FROM physical_files WHERE sha256_hash = $1 LIMIT 1;
-- name: CreatePhysicalFile :one
INSERT INTO physical_files (sha256_hash, size_bytes, storage_path) VALUES ($1, $2, $3) RETURNING *;
-- name: IncrementPhysicalFileRefCount :one
UPDATE physical_files SET reference_count = reference_count + 1 WHERE id = $1 RETURNING *;
-- name: CreateUserFile :one
INSERT INTO user_files (owner_id, physical_file_id, filename, mime_type, description, tags) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;
-- name: ListUserFiles :many
SELECT * FROM user_files WHERE owner_id = $1 ORDER BY upload_date DESC;
-- name: GetUserFileForDownload :one
SELECT uf.*, pf.storage_path FROM user_files uf JOIN physical_files pf ON uf.physical_file_id = pf.id WHERE uf.id = $1 AND uf.owner_id = $2;

-- name: GetFileOwnerAndPhysicalFile :one
-- CORRECTED: Added pf.storage_path to the SELECT and uf.owner_id to the WHERE clause.
SELECT uf.owner_id, pf.id as physical_file_id, pf.size_bytes, pf.storage_path
FROM user_files uf
JOIN physical_files pf ON uf.physical_file_id = pf.id
WHERE uf.id = $1 AND uf.owner_id = $2;

-- name: DeleteUserFile :exec
DELETE FROM user_files WHERE id = $1;
-- name: DecrementPhysicalFileRefCount :one
UPDATE physical_files SET reference_count = reference_count - 1 WHERE id = $1 RETURNING reference_count;
-- name: DeletePhysicalFile :exec
DELETE FROM physical_files WHERE id = $1;