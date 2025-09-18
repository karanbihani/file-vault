-- name: CreatePublicShareLink :one
INSERT INTO shares (user_file_id, share_token) VALUES ($1, $2) RETURNING *;

-- name: GetShareByToken :one
-- CORRECTED NAME: Changed from GetShareMetaByToken for clarity and consistency.
SELECT s.id, s.download_count, uf.filename, pf.storage_path, pf.size_bytes
FROM shares s
JOIN user_files uf ON s.user_file_id = uf.id
JOIN physical_files pf ON uf.physical_file_id = pf.id
WHERE s.share_token = $1 AND s.is_public = TRUE;

-- name: IncrementShareDownloadCount :exec
UPDATE shares SET download_count = download_count + 1 WHERE id = $1;