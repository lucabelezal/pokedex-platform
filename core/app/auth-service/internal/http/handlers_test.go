package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pokedex-platform/core/app/auth-service/internal/repository"
	"pokedex-platform/core/app/auth-service/internal/service"

	"github.com/golang-jwt/jwt/v5"
)

type stubAuthService struct {
	signupResult *service.AuthResult
	signupErr    error
	loginResult  *service.AuthResult
	loginErr     error
	refreshFn    func(ctx context.Context, token string) (*service.AuthResult, error)
	logoutFn     func(ctx context.Context, token string) error
}

func (s *stubAuthService) Signup(ctx context.Context, email, password string) (*service.AuthResult, error) {
	return s.signupResult, s.signupErr
}

func (s *stubAuthService) Login(ctx context.Context, email, password string) (*service.AuthResult, error) {
	return s.loginResult, s.loginErr
}

func (s *stubAuthService) Refresh(ctx context.Context, token string) (*service.AuthResult, error) {
	if s.refreshFn != nil {
		return s.refreshFn(ctx, token)
	}
	return nil, errors.New("not implemented")
}

func (s *stubAuthService) Logout(ctx context.Context, token string) error {
	if s.logoutFn != nil {
		return s.logoutFn(ctx, token)
	}
	return errors.New("not implemented")
}

func (s *stubAuthService) IsAccessTokenActive(ctx context.Context, tokenString string) (bool, error) {
	return true, nil
}

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

func TestSignup(t *testing.T) {
	tests := []struct {
		name       string
		signupRes  *service.AuthResult
		signupErr  error
		wantStatus int
	}{
		{
			name: "sucesso",
			signupRes: &service.AuthResult{
				UserID:       "user-1",
				Email:        "ash@kanto.dev",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "email duplicado",
			signupErr:  repository.ErrUserAlreadyExist,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "senha curta",
			signupErr:  service.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &stubAuthService{
				signupResult: tt.signupRes,
				signupErr:    tt.signupErr,
			}
			mux := NewMux(svc)

			body := strings.NewReader(`{"email":"ash@kanto.dev","password":"pikapika"}`)
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/signup", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusCreated {
				var result service.AuthResult
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Fatalf("falha ao decodificar: %v", err)
				}
				if result.UserID == "" {
					t.Error("UserID vazio")
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name       string
		loginRes   *service.AuthResult
		loginErr   error
		wantStatus int
	}{
		{
			name: "sucesso",
			loginRes: &service.AuthResult{
				UserID:       "user-1",
				Email:        "ash@kanto.dev",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "senha incorreta",
			loginErr:   service.ErrInvalidCredentials,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "entrada invalida",
			loginErr:   service.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &stubAuthService{
				loginResult: tt.loginRes,
				loginErr:    tt.loginErr,
			}
			mux := NewMux(svc)

			body := strings.NewReader(`{"email":"ash@kanto.dev","password":"pikapika"}`)
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var result service.AuthResult
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Fatalf("falha ao decodificar: %v", err)
				}
				if result.AccessToken == "" {
					t.Error("AccessToken vazio")
				}
			}
		})
	}
}

type stubPinger struct {
	err error
}

func (s stubPinger) Ping(ctx context.Context) error {
	return s.err
}

func TestRefreshHandler(t *testing.T) {
	tests := []struct {
		name       string
		refreshFn  func(ctx context.Context, token string) (*service.AuthResult, error)
		authHeader string
		wantStatus int
	}{
		{
			name: "sem token",
			refreshFn: func(ctx context.Context, token string) (*service.AuthResult, error) {
				return nil, nil
			},
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "token invalido",
			refreshFn: func(ctx context.Context, token string) (*service.AuthResult, error) {
				return nil, service.ErrInvalidToken
			},
			authHeader: "Bearer invalid",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "sucesso",
			refreshFn: func(ctx context.Context, token string) (*service.AuthResult, error) {
				return &service.AuthResult{UserID: "user-1", Email: "ash@kanto.dev", AccessToken: "new-access", RefreshToken: "new-refresh"}, nil
			},
			authHeader: "Bearer valid-token",
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &stubAuthService{refreshFn: tt.refreshFn}
			mux := NewMux(svc)
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	tests := []struct {
		name       string
		logoutFn   func(ctx context.Context, token string) error
		authHeader string
		wantStatus int
	}{
		{
			name:       "sem token",
			logoutFn:   func(ctx context.Context, token string) error { return nil },
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "token invalido",
			logoutFn:   func(ctx context.Context, token string) error { return service.ErrInvalidToken },
			authHeader: "Bearer invalid",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "erro interno",
			logoutFn:   func(ctx context.Context, token string) error { return errors.New("db down") },
			authHeader: "Bearer valid",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "sucesso",
			logoutFn:   func(ctx context.Context, token string) error { return nil },
			authHeader: "Bearer valid",
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &stubAuthService{logoutFn: tt.logoutFn}
			mux := NewMux(svc)
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHealthHandlerAuth(t *testing.T) {
	svc := &stubAuthService{}
	mux := NewMux(svc)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestReadyHandlerAuth(t *testing.T) {
	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
	}{
		{"banco pronto", nil, http.StatusOK},
		{"banco indisponivel", errors.New("connection refused"), http.StatusServiceUnavailable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ReadyHandler(stubPinger{err: tt.pingErr})
			req := httptest.NewRequest(http.MethodGet, "/ready", nil)
			w := httptest.NewRecorder()
			h(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestSignupPayloadInvalido(t *testing.T) {
	svc := &stubAuthService{}
	mux := NewMux(svc)
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/signup", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}
