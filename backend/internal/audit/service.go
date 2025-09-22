package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	Log(ctx context.Context, params LogParams)
}

type service struct {
	repo sqlc.Querier
}

func NewService(repo sqlc.Querier) Service {
	return &service{repo: repo}
}

type LogParams struct {
	UserID   int64
	Action   string
	TargetID uuid.UUID
	Details  map[string]interface{}
}

// The Log method runs in a separate goroutine so it never blocks the main request.
func (s *service) Log(ctx context.Context, params LogParams) {
	go func() {
		detailsJSON, err := json.Marshal(params.Details)
		if err != nil {
			log.Printf("CRITICAL: failed to marshal audit log details: %v", err)
			return
		}

		arg := sqlc.CreateAuditLogParams{
			UserID:   sql.NullInt64{Int64: params.UserID, Valid: true},
			Action:   sqlc.AuditAction(params.Action),
			TargetID: pgtype.UUID{Bytes: params.TargetID, Valid: true},
			Details:  detailsJSON,
		}

		// using a background context because the original request context might be cancelled.
		_, err = s.repo.CreateAuditLog(context.Background(), arg)
		if err != nil {
			log.Printf("CRITICAL: failed to create audit log: %v", err)
		}
	}()
}
