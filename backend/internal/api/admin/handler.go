package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apphandler"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes sets up the routing for the admin endpoints.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/audit-logs", apphandler.MakeHTTPHandler(h.GetAuditLogs))
	r.Get("/audit-logs/stats/activity-by-day", apphandler.MakeHTTPHandler(h.GetLogActivityStats))
}

// GetAuditLogs handles requests for the raw, paginated audit log feed.
func (h *Handler) GetAuditLogs(w http.ResponseWriter, r *http.Request) error {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	logs, err := h.service.ListAuditLogs(r.Context(), page, limit)
	if err != nil {
		return apierror.NewInternalServerError("Failed to retrieve audit logs")
	}

	return util.WriteJSON(w, http.StatusOK, logs)
}

// GetLogActivityStats handles requests for the daily activity graph data.
func (h *Handler) GetLogActivityStats(w http.ResponseWriter, r *http.Request) error {
	// Default to the last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// allow overriding with query parameters
	if sd := r.URL.Query().Get("start_date"); sd != "" {
		parsed, err := time.Parse("2006-01-02", sd)
		if err == nil {
			startDate = parsed
		}
	}
	if ed := r.URL.Query().Get("end_date"); ed != "" {
		parsed, err := time.Parse("2006-01-02", ed)
		if err == nil {
			endDate = parsed
		}
	}

	stats, err := h.service.GetLogActivityByDay(r.Context(), startDate, endDate)
	if err != nil {
		return apierror.NewInternalServerError("Failed to retrieve activity stats")
	}

	return util.WriteJSON(w, http.StatusOK, stats)
}
