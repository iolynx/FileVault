package apphandler

import (
	"errors"
	"log"
	"net/http"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
)

// AppHandler is a custom handler that returns an error.
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// MakeHTTPHandler converts an AppHandler into a standard http.HandlerFunc.
func MakeHTTPHandler(handler AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			var apiErr *apierror.APIError
			if errors.As(err, &apiErr) {
				util.WriteError(w, apiErr.StatusCode, apiErr.Error())
			} else {
				// log unexpected errors, return status 500.
				log.Printf("Unhandled error: %v", err)
				util.WriteError(w, http.StatusInternalServerError, "An unexpected error occurred")
			}
		}
	}
}
