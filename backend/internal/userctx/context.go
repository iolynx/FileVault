package userctx

import (
	"context"
	"strconv"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// SetUserID stores the userID (as string) in the request context.
func SetUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// GetUserID extracts the userID from context and converts it to int64.
// Returns (0, false) if missing or invalid.
func GetUserID(ctx context.Context) (int64, bool) {
	val := ctx.Value(userIDKey)
	if val == nil {
		return 0, false
	}

	idStr, ok := val.(string)
	if !ok {
		return 0, false
	}

	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, false
	}

	return idInt, true
}
