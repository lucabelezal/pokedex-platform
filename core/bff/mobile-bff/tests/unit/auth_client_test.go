package unit

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	httpclient "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/http"
	"pokedex-platform/core/bff/mobile-bff/internal/domain"

	"github.com/stretchr/testify/assert"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newTestAuthClient(fn roundTripFunc) *httpclient.AuthServiceClient {
	return httpclient.NewAuthServiceClientWithHTTPClient("http://auth-service.test", &http.Client{
		Transport: fn,
	})
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestAuthServiceClientSignupSuccess(t *testing.T) {
	client := newTestAuthClient(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/v1/auth/signup", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		return jsonResponse(http.StatusCreated, `{"accessToken":"token-123","refreshToken":"refresh-123","tokenType":"Bearer","expiresIn":900,"userId":"user-1","email":"ash@kanto.dev"}`), nil
	})

	session, err := client.Signup(context.Background(), "ash@kanto.dev", "pikachu123")

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "token-123", session.AccessToken)
	assert.Equal(t, "refresh-123", session.RefreshToken)
	assert.Equal(t, "user-1", session.UserID)
}

func TestAuthServiceClientErrorMapping(t *testing.T) {
	tests := []struct {
		name      string
		method    func(client *httpclient.AuthServiceClient) error
		transport roundTripFunc
		wantErr   error
	}{
		{
			name:    "login com 401 mapeia para ErrInvalidCredentials",
			wantErr: domain.ErrInvalidCredentials,
			transport: func(r *http.Request) (*http.Response, error) {
				return jsonResponse(http.StatusUnauthorized, `{"error":"credenciais invalidas"}`), nil
			},
			method: func(c *httpclient.AuthServiceClient) error {
				_, err := c.Login(context.Background(), "ash@kanto.dev", "wrong")
				return err
			},
		},
		{
			name:    "signup com 409 mapeia para ErrUserAlreadyExists",
			wantErr: domain.ErrUserAlreadyExists,
			transport: func(r *http.Request) (*http.Response, error) {
				return jsonResponse(http.StatusConflict, `{"error":"usuario ja existe"}`), nil
			},
			method: func(c *httpclient.AuthServiceClient) error {
				_, err := c.Signup(context.Background(), "ash@kanto.dev", "pikachu123")
				return err
			},
		},
		{
			name:    "refresh com 401 mapeia para ErrInvalidToken",
			wantErr: domain.ErrInvalidToken,
			transport: func(r *http.Request) (*http.Response, error) {
				return jsonResponse(http.StatusUnauthorized, `{"error":"token invalido"}`), nil
			},
			method: func(c *httpclient.AuthServiceClient) error {
				_, err := c.Refresh(context.Background(), "expired-token")
				return err
			},
		},
		{
			name:    "signup com 400 mapeia para ErrInvalidInput",
			wantErr: domain.ErrInvalidInput,
			transport: func(r *http.Request) (*http.Response, error) {
				return jsonResponse(http.StatusBadRequest, `{"error":"payload invalido"}`), nil
			},
			method: func(c *httpclient.AuthServiceClient) error {
				_, err := c.Signup(context.Background(), "", "")
				return err
			},
		},
		{
			name:    "logout com 503 mapeia para ErrAuthUnavailable",
			wantErr: domain.ErrAuthUnavailable,
			transport: func(r *http.Request) (*http.Response, error) {
				return jsonResponse(http.StatusServiceUnavailable, `{"error":"unavailable"}`), nil
			},
			method: func(c *httpclient.AuthServiceClient) error {
				return c.Logout(context.Background(), "token-123")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestAuthClient(tt.transport)
			err := tt.method(client)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
