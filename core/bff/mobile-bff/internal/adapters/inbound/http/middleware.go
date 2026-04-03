package httphandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDContextKey contextKey = "userID"
const UserEmailContextKey contextKey = "userEmail"

// rotas públicas que não passam por validação de token
var publicPaths = []string{
	"/health",
	"/api/v1/health",
	"/api/v1/auth/",
}

var authRateLimitedPaths = []string{
	"/api/v1/auth/login",
	"/api/v1/auth/signup",
	"/api/v1/auth/refresh",
	"/api/v1/auth/logout",
}

const (
	defaultAuthRateLimitRequests      = 20
	defaultAuthRateLimitWindowSeconds = 60
)

type rateLimitWindow struct {
	startedAt time.Time
	count     int
}

type authRateLimiter struct {
	m                 sync.Mutex
	entries           map[string]rateLimitWindow
	maxRequests       int
	window            time.Duration
	now               func() time.Time
	lastCleanupBucket int64
}

func isPublicPath(path string) bool {
	for _, prefix := range publicPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func isAuthRateLimitedPath(path string) bool {
	for _, candidate := range authRateLimitedPaths {
		if path == candidate {
			return true
		}
	}
	return false
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *statusResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// RequestLoggerMiddleware loga cada requisicao HTTP com metodo, path, status e duracao.
// Health checks sao ignorados para evitar ruido nos dashboards.
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" || r.URL.Path == "/api/v1/health" {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		rw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.InfoContext(r.Context(), "http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// AuthRateLimitMiddleware limita tentativas por IP nas rotas públicas de autenticação.
func AuthRateLimitMiddleware(next http.Handler) http.Handler {
	limiter := newAuthRateLimiter(getAuthRateLimitRequests(), getAuthRateLimitWindow())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuthRateLimitedPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		clientID := clientIdentifier(r)
		if !limiter.Allow(clientID) {
			slog.Warn("auth_audit", "action", "rate_limit", "outcome", "blocked", "client_ip", clientID, "path", r.URL.Path)
			RespondError(w, http.StatusTooManyRequests, "muitas tentativas, tente novamente em instantes", "TOO_MANY_REQUESTS")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware extrai e valida o token JWT da requisição.
// Rotas públicas (health e auth/*) passam sem validação.
func AuthMiddleware(validator inbound.TokenValidator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		tokenString, err := extractTokenFromRequest(r)
		if err != nil {
			RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
			return
		}

		var userID string
		var userEmail string
		if tokenString != "" {
			claims, err := parseAndValidateJWT(tokenString, getJWTSecret())
			if err != nil {
				RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
				return
			}

			active, activeErr := validator.IsTokenActive(r.Context(), tokenString)
			if activeErr != nil {
				RespondError(w, http.StatusServiceUnavailable, "auth service unavailable", "AUTH_UNAVAILABLE")
				return
			}
			if !active {
				RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
				return
			}

			userID = extractUserIDFromClaims(claims)
			if userID == "" {
				RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
				return
			}

			userEmail = extractUserEmailFromClaims(claims)
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		ctx = context.WithValue(ctx, UserEmailContextKey, userEmail)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader != "" {
		// Se Authorization foi enviado, ele tem prioridade sobre cookie.
		return parseBearerToken(authHeader)
	}

	cookie, err := r.Cookie("auth_token")
	if err == nil {
		tokenString := strings.TrimSpace(cookie.Value)
		if tokenString != "" {
			return tokenString, nil
		}
	}

	return "", nil
}

// RequireAuthMiddleware é um middleware que exige autenticação.
func RequireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())
		if userID == "" {
			http.Error(w, `{"error":"unauthorized","message":"autenticacao obrigatoria"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware adiciona headers CORS.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" && isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			if origin != "" && !isAllowedOrigin(origin) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
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

func extractUserEmailFromClaims(claims jwt.MapClaims) string {
	if email, ok := claims["email"].(string); ok {
		return strings.TrimSpace(email)
	}
	return ""
}

func getJWTSecret() string {
	return strings.TrimSpace(os.Getenv("JWT_SECRET"))
}

func getUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// SetUserID armazena um userID no contexto.
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}

func getUserEmailFromContext(ctx context.Context) string {
	userEmail, ok := ctx.Value(UserEmailContextKey).(string)
	if !ok {
		return ""
	}
	return userEmail
}

// SetUserEmail armazena um email no contexto.
func SetUserEmail(ctx context.Context, userEmail string) context.Context {
	return context.WithValue(ctx, UserEmailContextKey, userEmail)
}

// GetUserEmail retorna o email do usuário do contexto.
func GetUserEmail(ctx context.Context) string {
	return getUserEmailFromContext(ctx)
}

// GetUserID retorna o userID do contexto.
func GetUserID(ctx context.Context) string {
	return getUserIDFromContext(ctx)
}

func newAuthRateLimiter(maxRequests int, window time.Duration) *authRateLimiter {
	if maxRequests <= 0 {
		maxRequests = defaultAuthRateLimitRequests
	}
	if window <= 0 {
		window = time.Duration(defaultAuthRateLimitWindowSeconds) * time.Second
	}

	return &authRateLimiter{
		entries:     make(map[string]rateLimitWindow),
		maxRequests: maxRequests,
		window:      window,
		now:         time.Now,
	}
}

func (l *authRateLimiter) Allow(clientID string) bool {
	if clientID == "" {
		clientID = "unknown"
	}

	now := l.now()
	bucket := now.Unix() / int64(l.window.Seconds())

	l.m.Lock()
	defer l.m.Unlock()

	if bucket != l.lastCleanupBucket {
		l.cleanupExpired(now)
		l.lastCleanupBucket = bucket
	}

	entry, exists := l.entries[clientID]
	if !exists || now.Sub(entry.startedAt) >= l.window {
		l.entries[clientID] = rateLimitWindow{startedAt: now, count: 1}
		return true
	}

	if entry.count >= l.maxRequests {
		return false
	}

	entry.count++
	l.entries[clientID] = entry
	return true
}

func (l *authRateLimiter) cleanupExpired(now time.Time) {
	for clientID, entry := range l.entries {
		if now.Sub(entry.startedAt) >= l.window {
			delete(l.entries, clientID)
		}
	}
}

func clientIdentifier(r *http.Request) string {
	forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwardedFor != "" {
		first, _, _ := strings.Cut(forwardedFor, ",")
		if candidate := strings.TrimSpace(first); candidate != "" {
			return candidate
		}
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}

func getAuthRateLimitRequests() int {
	raw := strings.TrimSpace(os.Getenv("AUTH_RATE_LIMIT_REQUESTS"))
	if raw == "" {
		return defaultAuthRateLimitRequests
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return defaultAuthRateLimitRequests
	}

	return value
}

func getAuthRateLimitWindow() time.Duration {
	raw := strings.TrimSpace(os.Getenv("AUTH_RATE_LIMIT_WINDOW_SECONDS"))
	if raw == "" {
		return time.Duration(defaultAuthRateLimitWindowSeconds) * time.Second
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return time.Duration(defaultAuthRateLimitWindowSeconds) * time.Second
	}

	return time.Duration(value) * time.Second
}

func isAllowedOrigin(origin string) bool {
	allowedOrigins := getAllowedOrigins()
	if len(allowedOrigins) == 0 {
		return false
	}
	return slices.Contains(allowedOrigins, origin)
}

func getAllowedOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	if raw == "" {
		return []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:8000",
		}
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		origins = append(origins, origin)
	}

	return origins
}
