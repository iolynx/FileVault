package admin

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type AuditLogResponse struct {
	ID        int64                  `json:"id"`
	UserID    sql.NullInt64          `json:"user_id"`
	Action    string                 `json:"action"`
	TargetID  pgtype.UUID            `json:"target_id"`
	Details   map[string]interface{} `json:"details"` // The details are now a clean map
	CreatedAt time.Time              `json:"created_at"`
}
