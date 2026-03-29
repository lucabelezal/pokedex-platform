package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pokedex-platform/core/app/auth-service/internal/repository"
	"pokedex-platform/core/app/auth-service/internal/service"

	"github.com/golang-jwt/jwt/v5"
)

type stubAuthRepo struct {
	isAccessTokenRevokedFn func(ctx context.Context, jti string) (bool, error)
}

func (s *stubAuthRepo) CreateUser(ctx context.Context, email, passwordHash string) (*repository.User, error) {
	return nil, errors.New("not implemented")
}

func (s *stubAuthRepo) GetByEmail(ctx context.Context, email string) (*repository.User, error) {
	return nil, repository.ErrUserNotFound
}

func (s *stubAuthRepo) GetByID(ctx context.Context, userID string) (*repository.User, error) {
	return nil, repository.ErrUserNotFound
}

func (s *stubAuthRepo) StoreRefreshToken(ctx context.Context, userID, refreshToken string, expiresAt time.Time) error {
	return nil
}

func (s *stubAuthRepo) GetActiveRefreshSession(ctx context.Context, refreshToken string) (*repository.RefreshSession, error) {
	return nil, repository.ErrRefreshTokenNotFound
}

func (s *stubAuthRepo) RotateRefreshToken(ctx context.Context, currentToken, newToken, userID string, expiresAt time.Time) error {
	return repository.ErrRefreshTokenNotFound
}

func (s *stubAuthRepo) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	return repository.ErrRefreshTokenNotFound
}

func (s *stubAuthRepo) RevokeAccessToken(ctx context.Context, jti string, expiresAt time.Time) error {
	return nil
}

func (s *stubAuthRepo) IsAccessTokenRevoked(ctx context.Context, jti string) (bool, error) {
	if s.isAccessTokenRevokedFn != nil {
		return s.isAccessTokenRevokedFn(ctx, jti)
	}
	return false, nil
}

func TestIntrospectReturnsActiveTrueForValidToken(t *testing.T) {
	repo := &stubAuthRepo{
		isAccessTokenRevokedFn: func(ctx context.Context, jti string) (bool, error) {
			if jti != "jti-ativo" {
				t.Fatalf("jti inesperado: %s", jti)
			}
			return false, nil
		},
	}

	authService := service.NewAuthService(repo, "segredo", 15, 24)
	mux := NewMux(authService)
	token := mustSignToken(t, "segredo", "jti-ativo", time.Now().Add(10*time.Minute))

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/introspect", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", w.Code)
	}

	var resp introspectResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("falha ao decodificar resposta: %v", err)
	}
	if !resp.Active {
		t.Fatalf("token deveria estar ativo")
	}
}

func TestIntrospectReturnsActiveFalseForRevokedToken(t *testing.T) {
	repo := &stubAuthRepo{
		isAccessTokenRevokedFn: func(ctx context.Context, jti string) (bool, error) {
			if jti != "jti-revogado" {
				t.Fatalf("jti inesperado: %s", jti)
			}
			return true, nil
		},
	}

	authService := service.NewAuthService(repo, "segredo", 15, 24)
	mux := NewMux(authService)
	token := mustSignToken(t, "segredo", "jti-revogado", time.Now().Add(10*time.Minute))

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/introspect", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", w.Code)
	}

	var resp introspectResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("falha ao decodificar resposta: %v", err)
	}
	if resp.Active {
		t.Fatalf("token deveria estar inativo")
	}
}

func TestIntrospectWithoutAuthorizationReturnsInactive(t *testing.T) {
	authService := service.NewAuthService(&stubAuthRepo{}, "segredo", 15, 24)
	mux := NewMux(authService)

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/introspect", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", w.Code)
	}

	var resp introspectResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("falha ao decodificar resposta: %v", err)
	}
	if resp.Active {
		t.Fatalf("token sem autorizacao deveria ser inativo")
	}
}

func mustSignToken(t *testing.T, secret, jti string, expiresAt time.Time) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   "user-1",
		"email": "ash@kanto.dev",
		"iat":   time.Now().Unix(),
		"exp":   expiresAt.Unix(),
		"jti":   jti,
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("erro ao assinar token: %v", err)
	}

	return signed
}
