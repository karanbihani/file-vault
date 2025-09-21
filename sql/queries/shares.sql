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

-- name: ShareFileWithUser :exec
-- Creates a record in the junction table to share a file with a specific user.
INSERT INTO file_shares_to_users (
  user_file_id,
  shared_with_user_id
) VALUES (
  $1, $2
);

-- name: IsFileAlreadySharedWithUser :one
-- Checks if a share record already exists to prevent duplicates.
SELECT EXISTS(
  SELECT 1 FROM file_shares_to_users
  WHERE user_file_id = $1 AND shared_with_user_id = $2
);

-- name: UnshareFileWithUser :exec
-- Removes a specific user's access to a shared file.
DELETE FROM file_shares_to_users
WHERE user_file_id = $1 AND shared_with_user_id = $2;

-- name: DeletePublicShareLinksByFileID :exec
-- Removes ALL public share links associated with a specific file.
DELETE FROM shares
WHERE user_file_id = $1;