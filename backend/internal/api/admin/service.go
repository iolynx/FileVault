package admin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	ListAuditLogs(ctx context.Context, page, limit int) ([]AuditLogResponse, error)
	GetLogActivityByDay(ctx context.Context, startDate, endDate time.Time) ([]sqlc.GetAuditLogActivityByDayRow, error)
}

type service struct {
	repo sqlc.Querier
}

func NewService(repo sqlc.Querier) Service {
	return &service{repo: repo}
}

// ListAuditLogs handles the logic for paginating audit logs.
func (s *service) ListAuditLogs(ctx context.Context, page, limit int) ([]AuditLogResponse, error) {
	offset := (page - 1) * limit
	params := sqlc.ListAuditLogsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	rawLogs, err := s.repo.ListAuditLogs(ctx, params)
	if err != nil {
		return nil, err
	}

	// Transform the raw logs into our frontend-friendly DTO
	responseLogs := make([]AuditLogResponse, len(rawLogs))
	for i, rawLog := range rawLogs {
		var detailsMap map[string]interface{}
		// Unmarshal the raw JSONB bytes into our map
		if rawLog.Details != nil {
			err := json.Unmarshal(rawLog.Details, &detailsMap)
			if err != nil {
				// If there's an error, we can log it and/or set a default error message
				detailsMap = map[string]interface{}{"error": "failed to parse details"}
			}
		}

		responseLogs[i] = AuditLogResponse{
			ID:        rawLog.ID,
			UserID:    rawLog.UserID,
			Action:    string(rawLog.Action),
			TargetID:  rawLog.TargetID,
			Details:   detailsMap, // Assign the clean map
			CreatedAt: rawLog.CreatedAt.Time,
		}
	}

	return responseLogs, nil
}

// GetLogActivityByDay handles the logic for fetching daily stats.
func (s *service) GetLogActivityByDay(ctx context.Context, startDate, endDate time.Time) ([]sqlc.GetAuditLogActivityByDayRow, error) {
	params := sqlc.GetAuditLogActivityByDayParams{
		StartDate: pgtype.Timestamptz{Time: startDate, Valid: true},
		EndDate:   pgtype.Timestamptz{Time: endDate, Valid: true},
	}
	return s.repo.GetAuditLogActivityByDay(ctx, params)
}
