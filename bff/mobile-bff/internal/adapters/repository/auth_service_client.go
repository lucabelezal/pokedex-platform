package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
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
	return &AuthServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Signup chama o endpoint de signup do auth-service
func (c *AuthServiceClient) Signup(ctx context.Context, email, password string) (*AuthResponse, error) {
	if c.baseURL == "" {
		return nil, errors.New("auth service URL not configured")
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
		return nil, fmt.Errorf("signup request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read signup response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("signup failed: %s", errResp.Error)
		}
		return nil, fmt.Errorf("signup failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var authResp AuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse signup response: %w", err)
	}

	return &authResp, nil
}

// Login chama o endpoint de login do auth-service
func (c *AuthServiceClient) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	if c.baseURL == "" {
		return nil, errors.New("auth service URL not configured")
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
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read login response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("login failed: %s", errResp.Error)
		}
		return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var authResp AuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %w", err)
	}

	return &authResp, nil
}
