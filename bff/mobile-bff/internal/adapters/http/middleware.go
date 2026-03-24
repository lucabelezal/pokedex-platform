package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

// AuthMiddleware extrai e valida o token JWT da requisição
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		var userID string
		if authHeader != "" {
			tokenString, err := parseBearerToken(authHeader)
			if err != nil {
				RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
				return
			}

			claims, err := parseAndValidateJWT(tokenString, getJWTSecret())
			if err != nil {
				RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
				return
			}

			userID = extractUserIDFromClaims(claims)
			if userID == "" {
				RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth é um middleware que exige autenticação
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

// CORSMiddleware adiciona headers CORS
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

func parseBearerToken(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("authorization header nao usa bearer")
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenString == "" {
		return "", errors.New("token vazio")
	}

	return tokenString, nil
}

func parseAndValidateJWT(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de assinatura invalido: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("claims invalidas")
	}

	if expRaw, ok := claims["exp"]; ok {
		exp, ok := expRaw.(float64)
		if !ok {
			return nil, errors.New("exp invalido")
		}
		if time.Now().Unix() >= int64(exp) {
			return nil, errors.New("token expirado")
		}
	}

	return claims, nil
}

func extractUserIDFromClaims(claims jwt.MapClaims) string {
	if userID, ok := claims["user_id"].(string); ok && strings.TrimSpace(userID) != "" {
		return userID
	}

	if sub, ok := claims["sub"].(string); ok && strings.TrimSpace(sub) != "" {
		return sub
	}

	return ""
}

func getJWTSecret() string {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return "dev-secret"
	}
	return secret
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

func GetUserID(ctx context.Context) string {
	return getUserIDFromContext(ctx)
}
