package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpadapter "pokedex-platform/core/bff/mobile-bff/internal/adapters/inbound/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareWithToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user123",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("test-secret"))
	assert.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := httpadapter.GetUserID(r.Context())
		assert.Equal(t, "user123", userID)
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.AuthMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddlewareWithInvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.AuthMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token-invalido")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddlewareWithoutToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.AuthMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddlewareWithCookieToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user-cookie",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("test-secret"))
	assert.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := httpadapter.GetUserID(r.Context())
		assert.Equal(t, "user-cookie", userID)
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.AuthMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: tokenString})
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddlewareHeaderHasPriorityOverCookie(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user-cookie",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})
	validCookieToken, err := token.SignedString([]byte("test-secret"))
	assert.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.AuthMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token-invalido")
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: validCookieToken})
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCORSMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.CORSMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
}

func TestCORSMiddlewareOptions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called for OPTIONS")
	})

	middleware := httpadapter.CORSMiddleware(handler)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORSMiddlewareBlocksUnknownOriginOnPreflight(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.CORSMiddleware(handler)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthRateLimitMiddlewareAllowsWithinLimit(t *testing.T) {
	t.Setenv("AUTH_RATE_LIMIT_REQUESTS", "2")
	t.Setenv("AUTH_RATE_LIMIT_WINDOW_SECONDS", "60")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.AuthRateLimitMiddleware(handler)

	for range 2 {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
		req.RemoteAddr = "192.168.0.10:1234"
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestAuthRateLimitMiddlewareBlocksAfterLimit(t *testing.T) {
	t.Setenv("AUTH_RATE_LIMIT_REQUESTS", "1")
	t.Setenv("AUTH_RATE_LIMIT_WINDOW_SECONDS", "60")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.AuthRateLimitMiddleware(handler)

	firstReq := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	firstReq.RemoteAddr = "192.168.0.20:1234"
	firstResp := httptest.NewRecorder()
	middleware.ServeHTTP(firstResp, firstReq)
	assert.Equal(t, http.StatusOK, firstResp.Code)

	secondReq := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	secondReq.RemoteAddr = "192.168.0.20:1234"
	secondResp := httptest.NewRecorder()
	middleware.ServeHTTP(secondResp, secondReq)

	assert.Equal(t, http.StatusTooManyRequests, secondResp.Code)
	assert.Contains(t, secondResp.Body.String(), "TOO_MANY_REQUESTS")
}

func TestAuthRateLimitMiddlewareIgnoresNonAuthRoutes(t *testing.T) {
	t.Setenv("AUTH_RATE_LIMIT_REQUESTS", "1")
	t.Setenv("AUTH_RATE_LIMIT_WINDOW_SECONDS", "60")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := httpadapter.AuthRateLimitMiddleware(handler)

	for range 3 {
		req := httptest.NewRequest("GET", "/api/v1/pokemons", nil)
		req.RemoteAddr = "192.168.0.30:1234"
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}
}
