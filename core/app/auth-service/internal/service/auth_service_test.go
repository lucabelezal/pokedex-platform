package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"pokedex-platform/core/app/auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type mockAuthRepository struct {
	createUserFn           func(ctx context.Context, email, passwordHash string) (*repository.User, error)
	getByEmailFn           func(ctx context.Context, email string) (*repository.User, error)
	getByIDFn              func(ctx context.Context, userID string) (*repository.User, error)
	storeRefreshTokenFn    func(ctx context.Context, userID, refreshToken string, expiresAt time.Time) error
	getActiveRefreshFn     func(ctx context.Context, refreshToken string) (*repository.RefreshSession, error)
	rotateRefreshTokenFn   func(ctx context.Context, currentToken, newToken, userID string, expiresAt time.Time) error
	revokeRefreshTokenFn   func(ctx context.Context, refreshToken string) error
	revokeAccessTokenFn    func(ctx context.Context, jti string, expiresAt time.Time) error
	isAccessTokenRevokedFn func(ctx context.Context, jti string) (bool, error)
}

func (m *mockAuthRepository) CreateUser(ctx context.Context, email, passwordHash string) (*repository.User, error) {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, email, passwordHash)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuthRepository) GetByEmail(ctx context.Context, email string) (*repository.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockAuthRepository) GetByID(ctx context.Context, userID string) (*repository.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, userID)
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockAuthRepository) StoreRefreshToken(ctx context.Context, userID, refreshToken string, expiresAt time.Time) error {
	if m.storeRefreshTokenFn != nil {
		return m.storeRefreshTokenFn(ctx, userID, refreshToken, expiresAt)
	}
	return nil
}

func (m *mockAuthRepository) GetActiveRefreshSession(ctx context.Context, refreshToken string) (*repository.RefreshSession, error) {
	if m.getActiveRefreshFn != nil {
		return m.getActiveRefreshFn(ctx, refreshToken)
	}
	return nil, repository.ErrRefreshTokenNotFound
}

func (m *mockAuthRepository) RotateRefreshToken(ctx context.Context, currentToken, newToken, userID string, expiresAt time.Time) error {
	if m.rotateRefreshTokenFn != nil {
		return m.rotateRefreshTokenFn(ctx, currentToken, newToken, userID, expiresAt)
	}
	return nil
}

func (m *mockAuthRepository) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if m.revokeRefreshTokenFn != nil {
		return m.revokeRefreshTokenFn(ctx, refreshToken)
	}
	return repository.ErrRefreshTokenNotFound
}

func (m *mockAuthRepository) RevokeAccessToken(ctx context.Context, jti string, expiresAt time.Time) error {
	if m.revokeAccessTokenFn != nil {
		return m.revokeAccessTokenFn(ctx, jti, expiresAt)
	}
	return nil
}

func (m *mockAuthRepository) IsAccessTokenRevoked(ctx context.Context, jti string) (bool, error) {
	if m.isAccessTokenRevokedFn != nil {
		return m.isAccessTokenRevokedFn(ctx, jti)
	}
	return false, nil
}

func TestRefreshRotatesRefreshToken(t *testing.T) {
	const oldRefreshToken = "refresh-antigo"

	rotateCalled := false
	repo := &mockAuthRepository{
		getActiveRefreshFn: func(ctx context.Context, refreshToken string) (*repository.RefreshSession, error) {
			if refreshToken != oldRefreshToken {
				t.Fatalf("refresh token inesperado: %s", refreshToken)
			}
			return &repository.RefreshSession{UserID: "user-1", ExpiresAt: time.Now().Add(1 * time.Hour)}, nil
		},
		getByIDFn: func(ctx context.Context, userID string) (*repository.User, error) {
			if userID != "user-1" {
				t.Fatalf("user id inesperado: %s", userID)
			}
			return &repository.User{ID: "user-1", Email: "ash@kanto.dev"}, nil
		},
		rotateRefreshTokenFn: func(ctx context.Context, currentToken, newToken, userID string, expiresAt time.Time) error {
			rotateCalled = true
			if currentToken != oldRefreshToken {
				t.Fatalf("token atual inesperado: %s", currentToken)
			}
			if newToken == "" || newToken == oldRefreshToken {
				t.Fatalf("novo refresh token invalido")
			}
			if userID != "user-1" {
				t.Fatalf("user id inesperado na rotacao: %s", userID)
			}
			return nil
		},
	}

	svc := NewAuthService(repo, "segredo-teste", 15, 24)
	result, err := svc.Refresh(context.Background(), oldRefreshToken)
	if err != nil {
		t.Fatalf("erro inesperado no refresh: %v", err)
	}
	if !rotateCalled {
		t.Fatalf("rotacao de refresh token nao foi executada")
	}
	if result.AccessToken == "" {
		t.Fatalf("access token nao retornado")
	}
	if result.RefreshToken == "" || result.RefreshToken == oldRefreshToken {
		t.Fatalf("refresh token retornado invalido")
	}
}

func TestLogoutRevogaAccessTokenQuandoRefreshNaoExiste(t *testing.T) {
	revokedJTI := ""
	repo := &mockAuthRepository{
		revokeRefreshTokenFn: func(ctx context.Context, refreshToken string) error {
			return repository.ErrRefreshTokenNotFound
		},
		revokeAccessTokenFn: func(ctx context.Context, jti string, expiresAt time.Time) error {
			revokedJTI = jti
			if expiresAt.Before(time.Now().UTC()) {
				t.Fatalf("expiracao invalida para revogacao de access token")
			}
			return nil
		},
	}

	token := mustSignAccessToken(t, "segredo-teste", "jti-logout", time.Now().Add(15*time.Minute))
	svc := NewAuthService(repo, "segredo-teste", 15, 24)

	if err := svc.Logout(context.Background(), token); err != nil {
		t.Fatalf("logout retornou erro: %v", err)
	}
	if revokedJTI != "jti-logout" {
		t.Fatalf("jti revogado inesperado: %s", revokedJTI)
	}
}

func TestIsAccessTokenActiveRetornaFalseQuandoRevogado(t *testing.T) {
	repo := &mockAuthRepository{
		isAccessTokenRevokedFn: func(ctx context.Context, jti string) (bool, error) {
			if jti != "jti-revogado" {
				t.Fatalf("jti inesperado em consulta de revogacao: %s", jti)
			}
			return true, nil
		},
	}

	token := mustSignAccessToken(t, "segredo-teste", "jti-revogado", time.Now().Add(15*time.Minute))
	svc := NewAuthService(repo, "segredo-teste", 15, 24)

	active, err := svc.IsAccessTokenActive(context.Background(), token)
	if err != nil {
		t.Fatalf("erro inesperado na introspeccao: %v", err)
	}
	if active {
		t.Fatalf("token deveria estar inativo apos revogacao")
	}
}

func TestLogoutComTokenInvalidoRetornaErro(t *testing.T) {
	repo := &mockAuthRepository{
		revokeRefreshTokenFn: func(ctx context.Context, refreshToken string) error {
			return repository.ErrRefreshTokenNotFound
		},
	}

	svc := NewAuthService(repo, "segredo-teste", 15, 24)
	err := svc.Logout(context.Background(), "token-invalido")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("esperava ErrInvalidToken, obteve: %v", err)
	}
}

func mustSignAccessToken(t *testing.T, secret, jti string, expiresAt time.Time) string {
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
		t.Fatalf("falha ao assinar token de teste: %v", err)
	}

	return signed
}
