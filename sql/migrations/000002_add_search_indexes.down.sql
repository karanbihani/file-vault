-- This migration rolls back the indexes created in the corresponding .up.sql file.
DROP INDEX IF EXISTS idx_user_files_filename;
DROP INDEX IF EXISTS idx_user_files_mime_type;
DROP INDEX IF EXISTS idx_user_files_upload_date;
DROP INDEX IF EXISTS idx_user_files_tags;
DROP INDEX IF EXISTS idx_physical_files_size_bytes;