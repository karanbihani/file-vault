-- This migration adds indexes to optimize search and filtering performance.

-- A B-Tree index on the filename column using the 'text_pattern_ops' operator class,
-- which is optimized for ILIKE pattern matching queries (e.g., 'foo%').
CREATE INDEX idx_user_files_filename ON user_files (filename text_pattern_ops);

-- Standard B-Tree indexes for exact matches and range queries.
CREATE INDEX idx_user_files_mime_type ON user_files (mime_type);
CREATE INDEX idx_user_files_upload_date ON user_files (upload_date);

-- A GIN index on the tags array. This is highly efficient for queries
-- that check for the presence of elements within the array (e.g., using the @> operator).
CREATE INDEX idx_user_files_tags ON user_files USING GIN (tags);

-- We also need an index on the physical_files table for size-based filtering.
CREATE INDEX idx_physical_files_size_bytes ON physical_files (size_bytes);