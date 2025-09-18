package middleware

import (
	"log"
	"net/http"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Auth Header & Process
			cookie, err := r.Cookie("jwt")
			if err != nil {
				util.WriteError(w, http.StatusUnauthorized, "Missing JWT Cookie")
				return
			}
			tokenString := cookie.Value

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrInvalidKeyType
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				util.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Print("invalid claims")
				util.WriteError(w, http.StatusUnauthorized, "invalid claims")
				return
			}

			userID, ok := claims["user_id"].(float64)
			if !ok {
				log.Print("invalid user_id claim")
				util.WriteError(w, http.StatusUnauthorized, "invalid user_id claim")
				return
			}

			// converting float64 to int64 to store in user context
			ctx := userctx.SetUserID(r.Context(), int64(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
