-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    user_id,
    action,
    target_id,
    details
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListAuditLogs :many
SELECT * FROM audit_logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAuditLogActivityByDay :many
SELECT
    DATE_TRUNC('day', created_at)::date AS activity_day,
    COUNT(*) AS event_count
FROM audit_logs
WHERE created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date')
GROUP BY activity_day
ORDER BY activity_day;
