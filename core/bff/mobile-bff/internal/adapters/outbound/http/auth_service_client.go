package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)

// AuthServiceClient fornece cliente HTTP para comunicação com auth-service.
type AuthServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	UserID       string `json:"userId"`
	Email        string `json:"email"`
}

// NewAuthServiceClient cria um novo cliente de auth-service.
func NewAuthServiceClient(baseURL string) *AuthServiceClient {
	return NewAuthServiceClientWithHTTPClient(baseURL, &http.Client{
		Timeout: 10 * time.Second,
	})
}

// NewAuthServiceClientWithHTTPClient cria um cliente com HTTP client injetado.
func NewAuthServiceClientWithHTTPClient(baseURL string, httpClient *http.Client) *AuthServiceClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &AuthServiceClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// Signup chama o endpoint de signup do auth-service.
func (c *AuthServiceClient) Signup(ctx context.Context, email, password string) (*domain.AuthSession, error) {
	if c.baseURL == "" {
		return nil, domain.ErrAuthUnavailable
	}

	body, err := json.Marshal(signupRequest{Email: email, Password: password})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signup request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/auth/signup", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create signup request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, domain.ErrAuthUnavailable
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read signup response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, mapAuthError(resp.StatusCode, "signup")
	}

	var ar authResponse
	if err := json.Unmarshal(respBody, &ar); err != nil {
		return nil, fmt.Errorf("failed to parse signup response: %w", err)
	}

	return toAuthSession(ar), nil
}

// Login chama o endpoint de login do auth-service.
func (c *AuthServiceClient) Login(ctx context.Context, email, password string) (*domain.AuthSession, error) {
	if c.baseURL == "" {
		return nil, domain.ErrAuthUnavailable
	}

	body, err := json.Marshal(loginRequest{Email: email, Password: password})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/auth/login", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, domain.ErrAuthUnavailable
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read login response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, mapAuthError(resp.StatusCode, "login")
	}

	var ar authResponse
	if err := json.Unmarshal(respBody, &ar); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %w", err)
	}

	return toAuthSession(ar), nil
}

// Refresh chama o endpoint de refresh do auth-service.
func (c *AuthServiceClient) Refresh(ctx context.Context, token string) (*domain.AuthSession, error) {
	if c.baseURL == "" {
		return nil, domain.ErrAuthUnavailable
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/auth/refresh", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, domain.ErrAuthUnavailable
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, mapAuthError(resp.StatusCode, "refresh")
	}

	var ar authResponse
	if err := json.Unmarshal(respBody, &ar); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	return toAuthSession(ar), nil
}

// Logout chama o endpoint de logout do auth-service.
func (c *AuthServiceClient) Logout(ctx context.Context, token string) error {
	if c.baseURL == "" {
		return domain.ErrAuthUnavailable
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/auth/logout", nil)
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return domain.ErrAuthUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return mapAuthError(resp.StatusCode, "logout")
	}

	return nil
}

func toAuthSession(ar authResponse) *domain.AuthSession {
	return &domain.AuthSession{
		AccessToken:  ar.AccessToken,
		RefreshToken: ar.RefreshToken,
		TokenType:    ar.TokenType,
		ExpiresIn:    ar.ExpiresIn,
		UserID:       ar.UserID,
		Email:        ar.Email,
	}
}

func mapAuthError(statusCode int, operation string) error {
	switch statusCode {
	case http.StatusBadRequest:
		return domain.ErrInvalidInput
	case http.StatusUnauthorized:
		if operation == "login" {
			return domain.ErrInvalidCredentials
		}
		return domain.ErrInvalidToken
	case http.StatusConflict:
		return domain.ErrUserAlreadyExists
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return domain.ErrAuthUnavailable
	default:
		if statusCode >= 500 {
			return domain.ErrAuthUnavailable
		}
		return fmt.Errorf("auth service returned status %d", statusCode)
	}
}

var _ outbound.AuthProvider = (*AuthServiceClient)(nil)
