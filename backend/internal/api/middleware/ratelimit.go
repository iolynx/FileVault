package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/redis/go-redis/v9"
)

// RateLimiter creates a middleware that enforces rate limits per user.
// It uses the environment variables `API_RATE_LIMIT=2` and  `API_RATE_LIMIT_WINDOW_SECONDS=1`
// to enforce configurable rate limiting, and returns a 429 StatusTooManyRequests upon
// exceeding the allowed number of requests within the configured time window.
func RateLimiter(redisClient *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get UserID from context.
			userID, ok := userctx.GetUserID(ctx)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			// Generate the Redis key for the user
			key := fmt.Sprintf("rate_limit:%d", userID)

			pipe := redisClient.TxPipeline()
			count := pipe.Incr(ctx, key)
			// Set the key to expire after the window duration
			pipe.Expire(ctx, key, window)

			_, err := pipe.Exec(ctx)
			if err != nil {
				log.Printf("Error executing Redis pipeline for rate limiting: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			// Check if the count exceeds the limit, if it does, return 429
			if count.Val() > int64(limit) {
				errResponse := apierror.New(http.StatusTooManyRequests, "Rate limit exceeded")
				util.WriteError(w, errResponse.StatusCode, errResponse.Message)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
