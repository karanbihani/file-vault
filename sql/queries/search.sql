-- name: SearchFiles :many
-- Performs a comprehensive search and filter operation on user files.
-- This query is optimized with indexes and uses sqlc.narg() for optional parameters.
SELECT
    uf.id,
    uf.filename,
    uf.mime_type,
    uf.upload_date,
    pf.size_bytes,
    u.email as owner_email
FROM
    user_files uf
JOIN
    users u ON uf.owner_id = u.id
JOIN
    physical_files pf ON uf.physical_file_id = pf.id
WHERE
    (
        @is_admin::boolean OR
        uf.owner_id = @requesting_user_id::bigint OR
        EXISTS (
            SELECT 1 FROM file_shares_to_users fstu
            WHERE fstu.user_file_id = uf.id AND fstu.shared_with_user_id = @requesting_user_id::bigint
        )
    )
AND
    (uf.filename ILIKE '%' || sqlc.narg('filename') || '%' OR sqlc.narg('filename') IS NULL)
AND
    (uf.mime_type = sqlc.narg('mime_type') OR sqlc.narg('mime_type') IS NULL)
AND
    (pf.size_bytes >= sqlc.narg('min_size') OR sqlc.narg('min_size') IS NULL)
AND
    (pf.size_bytes <= sqlc.narg('max_size') OR sqlc.narg('max_size') IS NULL)
AND
    (uf.upload_date >= sqlc.narg('start_date') OR sqlc.narg('start_date') IS NULL)
AND
    (uf.upload_date <= sqlc.narg('end_date') OR sqlc.narg('end_date') IS NULL)
AND
    -- The @> operator checks if the tags array contains all elements from the input array.
    -- This is efficiently powered by our GIN index.
    (uf.tags @> sqlc.narg('tags')::text[] OR sqlc.narg('tags') IS NULL)
AND
    -- Filter by a specific uploader's email if provided.
    (u.email = sqlc.narg('uploader_email') OR sqlc.narg('uploader_email') IS NULL)
ORDER BY
    uf.upload_date DESC;

    