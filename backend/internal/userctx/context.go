package userctx

import (
	"context"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

// SetUserID stores the userID (as int64) in the request context.
func SetUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// GetUserID extracts the userID from context
// Returns (0, false) if missing or invalid.
func GetUserID(ctx context.Context) (int64, bool) {
	val := ctx.Value(userIDKey)
	if val == nil {
		return 0, false
	}

	idInt, ok := val.(int64)
	if !ok {
		return 0, false
	}

	return idInt, true
}
