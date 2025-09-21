-- name: CreateAuditLog :exec
-- Inserts a new audit log entry.
INSERT INTO audit_logs (user_id, action, details)
VALUES ($1, $2, $3);

-- name: ListAuditLogs :many
-- For admin use: retrieves all audit log entries, newest first.
SELECT * FROM audit_logs
ORDER BY timestamp DESC;