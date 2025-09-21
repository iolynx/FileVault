package middleware

import (
	"net/http"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
)

// AdminMiddleware checks if the authenticated user has the 'admin' role.
// It runs after the AuthMiddleware.
func AdminMiddleware(repo sqlc.Querier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// This wont happen as AuthMiddleware runs before this,
			// but we are checking just in case.
			userID, ok := userctx.GetUserID(ctx)
			if !ok {
				errResponse := apierror.NewUnauthorizedError()
				util.WriteError(w, errResponse.StatusCode, errResponse.Message)
				return
			}

			// Get the user object from the db
			user, err := repo.GetUserByID(ctx, userID)
			if err != nil {
				errResponse := apierror.NewInternalServerError("could not retrieve user")
				util.WriteError(w, errResponse.StatusCode, errResponse.Message)
				return
			}

			// Check the user's role
			if user.Role != "admin" {
				// If not an admin, return a 403 Forbidden error
				errResponse := apierror.NewForbiddenError()
				util.WriteError(w, errResponse.StatusCode, errResponse.Message)
				return
			}

			// User is an admin, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
