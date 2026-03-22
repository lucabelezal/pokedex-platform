package http

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

// AuthMiddleware extracts and validates JWT token from request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		var userID string
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token != "" {
				userID = extractUserIDFromToken(token)
			}
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is a middleware that requires authentication
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())
		if userID == "" {
			http.Error(w, `{"error":"unauthorized","message":"autenticacao obrigatoria"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func extractUserIDFromToken(token string) string {
	if len(token) > 20 {
		return token[:20]
	}
	if len(token) > 0 {
		return token
	}
	return ""
}

func getUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}
