-- name: ListAllFiles :many
-- For admin use: retrieves all files with uploader's email.
SELECT
    uf.id,
    uf.filename,
    uf.mime_type,
    uf.upload_date,
    pf.size_bytes,
    u.email as owner_email
FROM user_files uf
JOIN users u ON uf.owner_id = u.id
JOIN physical_files pf ON uf.physical_file_id = pf.id
ORDER BY uf.upload_date DESC;

-- name: GetSystemStats :one
-- For admin use: retrieves system-wide aggregate statistics.
SELECT
    (SELECT COUNT(*) FROM users)::bigint AS total_users,
    (SELECT COUNT(*) FROM user_files)::bigint AS total_files,
    (SELECT COALESCE(SUM(size_bytes), 0) FROM physical_files)::bigint AS total_storage_used,
    (SELECT COALESCE(SUM(download_count), 0) FROM shares)::bigint AS total_downloads;

-- name: GetFileMetadataByID :one
-- For admin use: retrieves file metadata without any ownership checks.
SELECT
    uf.filename,
    pf.storage_path,
    pf.size_bytes
FROM user_files uf
JOIN physical_files pf ON uf.physical_file_id = pf.id
WHERE uf.id = $1;