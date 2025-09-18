-- name: GetUserStats :one
-- Retrieves storage statistics for a single user.
-- It gets the pre-calculated deduplicated usage from the users table
-- and calculates the original total size by summing up the sizes of all files owned by the user.
SELECT
    u.storage_used_bytes AS deduplicated_usage,
    COALESCE(SUM(pf.size_bytes), 0)::bigint AS original_usage
FROM users u
LEFT JOIN user_files uf ON u.id = uf.owner_id
LEFT JOIN physical_files pf ON uf.physical_file_id = pf.id
WHERE u.id = $1
GROUP BY u.id;

