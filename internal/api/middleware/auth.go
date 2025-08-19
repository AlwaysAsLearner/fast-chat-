package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/AlwaysAsLearner/fast-chat/backend/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %s", err), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
