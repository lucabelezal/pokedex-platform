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

func TestAuthServiceClientLoginMapsUnauthorizedToInvalidCredentials(t *testing.T) {
	client := newTestAuthClient(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/v1/auth/login", r.URL.Path)
		return jsonResponse(http.StatusUnauthorized, `{"error":"credenciais invalidas"}`), nil
	})

	session, err := client.Login(context.Background(), "ash@kanto.dev", "wrong")

	assert.Nil(t, session)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
}

func TestAuthServiceClientSignupMapsConflictToUserAlreadyExists(t *testing.T) {
	client := newTestAuthClient(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/v1/auth/signup", r.URL.Path)
		return jsonResponse(http.StatusConflict, `{"error":"usuario ja existe"}`), nil
	})

	session, err := client.Signup(context.Background(), "ash@kanto.dev", "pikachu123")

	assert.Nil(t, session)
	assert.Equal(t, domain.ErrUserAlreadyExists, err)
}

func TestAuthServiceClientRefreshMapsUnauthorizedToInvalidToken(t *testing.T) {
	client := newTestAuthClient(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/v1/auth/refresh", r.URL.Path)
		assert.Equal(t, "Bearer expired-token", r.Header.Get("Authorization"))
		return jsonResponse(http.StatusUnauthorized, `{"error":"token invalido"}`), nil
	})

	session, err := client.Refresh(context.Background(), "expired-token")

	assert.Nil(t, session)
	assert.Equal(t, domain.ErrInvalidToken, err)
}

func TestAuthServiceClientSignupMapsBadRequestToInvalidInput(t *testing.T) {
	client := newTestAuthClient(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusBadRequest, `{"error":"payload invalido"}`), nil
	})

	session, err := client.Signup(context.Background(), "", "")

	assert.Nil(t, session)
	assert.Equal(t, domain.ErrInvalidInput, err)
}

func TestAuthServiceClientLogoutMapsServiceUnavailableToAuthUnavailable(t *testing.T) {
	client := newTestAuthClient(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/v1/auth/logout", r.URL.Path)
		return jsonResponse(http.StatusServiceUnavailable, `{"error":"unavailable"}`), nil
	})

	err := client.Logout(context.Background(), "token-123")

	assert.Equal(t, domain.ErrAuthUnavailable, err)
}
