package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	"pokedex-platform/core/bff/mobile-bff/internal/ports"
)

// AuthServiceClient fornece cliente HTTP para comunicação com auth-service
type AuthServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// SignupRequest é o corpo da requisição de signup
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest é o corpo da requisição de login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse é a resposta do auth-service
type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
	UserID      string `json:"userId"`
	Email       string `json:"email"`
}

// ErrorResponse é a resposta de erro do auth-service
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewAuthServiceClient cria um novo cliente de auth-service
func NewAuthServiceClient(baseURL string) *AuthServiceClient {
	return NewAuthServiceClientWithHTTPClient(baseURL, &http.Client{
		Timeout: 10 * time.Second,
	})
}

// NewAuthServiceClientWithHTTPClient creates a client with an injected HTTP client.
func NewAuthServiceClientWithHTTPClient(baseURL string, httpClient *http.Client) *AuthServiceClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &AuthServiceClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// Signup chama o endpoint de signup do auth-service
func (c *AuthServiceClient) Signup(ctx context.Context, email, password string) (*ports.AuthSession, error) {
	if c.baseURL == "" {
		return nil, domain.ErrAuthUnavailable
	}

	req := SignupRequest{
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signup request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/auth/signup", bytes.NewReader(body))
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

	var authResp AuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse signup response: %w", err)
	}

	return &ports.AuthSession{
		AccessToken: authResp.AccessToken,
		TokenType:   authResp.TokenType,
		ExpiresIn:   authResp.ExpiresIn,
		UserID:      authResp.UserID,
		Email:       authResp.Email,
	}, nil
}

// Login chama o endpoint de login do auth-service
func (c *AuthServiceClient) Login(ctx context.Context, email, password string) (*ports.AuthSession, error) {
	if c.baseURL == "" {
		return nil, domain.ErrAuthUnavailable
	}

	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/auth/login", bytes.NewReader(body))
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

	var authResp AuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %w", err)
	}

	return &ports.AuthSession{
		AccessToken: authResp.AccessToken,
		TokenType:   authResp.TokenType,
		ExpiresIn:   authResp.ExpiresIn,
		UserID:      authResp.UserID,
		Email:       authResp.Email,
	}, nil
}

// Refresh chama o endpoint de refresh do auth-service
func (c *AuthServiceClient) Refresh(ctx context.Context, token string) (*ports.AuthSession, error) {
	if c.baseURL == "" {
		return nil, domain.ErrAuthUnavailable
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/auth/refresh", nil)
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

	var authResp AuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	return &ports.AuthSession{
		AccessToken: authResp.AccessToken,
		TokenType:   authResp.TokenType,
		ExpiresIn:   authResp.ExpiresIn,
		UserID:      authResp.UserID,
		Email:       authResp.Email,
	}, nil
}

// Logout chama o endpoint de logout do auth-service
func (c *AuthServiceClient) Logout(ctx context.Context, token string) error {
	if c.baseURL == "" {
		return domain.ErrAuthUnavailable
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/auth/logout", nil)
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

var _ ports.AuthProvider = (*AuthServiceClient)(nil)
