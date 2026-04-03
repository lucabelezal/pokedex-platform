package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

type authOperation string

const (
	authOperationSignup  authOperation = "signup"
	authOperationLogin   authOperation = "login"
	authOperationRefresh authOperation = "refresh"
	authOperationLogout  authOperation = "logout"
)

const maxAuthPayloadBytes int64 = 8 * 1024

// Signup gerencia registro de usuário.
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "metodo nao permitido", "METHOD_NOT_ALLOWED")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := decodeAuthJSONBody(w, r, &req); err != nil {
		RespondError(w, http.StatusBadRequest, "email e password obrigatorios", "INVALID_REQUEST")
		return
	}

	if req.Email == "" || req.Password == "" {
		RespondError(w, http.StatusBadRequest, "email e password obrigatorios", "INVALID_REQUEST")
		return
	}

	if h.authUseCase == nil {
		h.respondAuthError(w, domain.ErrAuthUnavailable, authOperationSignup)
		return
	}

	authResp, err := h.authUseCase.Signup(ctx, req.Email, req.Password)
	if err != nil {
		h.respondAuthError(w, err, authOperationSignup)
		return
	}

	http.SetCookie(w, buildAuthCookie(r, authResp.AccessToken, authResp.ExpiresIn))
	if authResp.RefreshToken != "" {
		http.SetCookie(w, buildRefreshCookie(r, authResp.RefreshToken))
	}

	RespondJSON(w, http.StatusCreated, struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		UserID      string `json:"user_id"`
		Email       string `json:"email"`
	}{
		AccessToken: authResp.AccessToken,
		TokenType:   authResp.TokenType,
		ExpiresIn:   authResp.ExpiresIn,
		UserID:      authResp.UserID,
		Email:       authResp.Email,
	})
}

// Login gerencia autenticação de usuário.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "metodo nao permitido", "METHOD_NOT_ALLOWED")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := decodeAuthJSONBody(w, r, &req); err != nil {
		RespondError(w, http.StatusBadRequest, "email e password obrigatorios", "INVALID_REQUEST")
		return
	}

	if req.Email == "" || req.Password == "" {
		RespondError(w, http.StatusBadRequest, "email e password obrigatorios", "INVALID_REQUEST")
		return
	}

	if h.authUseCase == nil {
		h.respondAuthError(w, domain.ErrAuthUnavailable, authOperationLogin)
		return
	}

	authResp, err := h.authUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.respondAuthError(w, err, authOperationLogin)
		return
	}

	http.SetCookie(w, buildAuthCookie(r, authResp.AccessToken, authResp.ExpiresIn))
	if authResp.RefreshToken != "" {
		http.SetCookie(w, buildRefreshCookie(r, authResp.RefreshToken))
	}

	RespondJSON(w, http.StatusOK, struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		UserID      string `json:"user_id"`
		Email       string `json:"email"`
	}{
		AccessToken: authResp.AccessToken,
		TokenType:   authResp.TokenType,
		ExpiresIn:   authResp.ExpiresIn,
		UserID:      authResp.UserID,
		Email:       authResp.Email,
	})
}

// Refresh renova o token de acesso do usuário autenticado.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "metodo nao permitido", "METHOD_NOT_ALLOWED")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tokenString, err := extractRefreshTokenFromRequest(r)
	if err != nil || tokenString == "" {
		RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
		return
	}

	if h.authUseCase == nil {
		h.respondAuthError(w, domain.ErrAuthUnavailable, authOperationRefresh)
		return
	}

	authResp, err := h.authUseCase.Refresh(ctx, tokenString)
	if err != nil {
		h.respondAuthError(w, err, authOperationRefresh)
		return
	}

	http.SetCookie(w, buildAuthCookie(r, authResp.AccessToken, authResp.ExpiresIn))
	if authResp.RefreshToken != "" {
		http.SetCookie(w, buildRefreshCookie(r, authResp.RefreshToken))
	}

	RespondJSON(w, http.StatusOK, struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		UserID      string `json:"user_id"`
		Email       string `json:"email"`
	}{
		AccessToken: authResp.AccessToken,
		TokenType:   authResp.TokenType,
		ExpiresIn:   authResp.ExpiresIn,
		UserID:      authResp.UserID,
		Email:       authResp.Email,
	})
}

// Logout encerra a sessão e remove os cookies de autenticação.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "metodo nao permitido", "METHOD_NOT_ALLOWED")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	accessToken, _ := extractTokenFromRequest(r)
	refreshToken, _ := extractRefreshTokenFromRequest(r)

	if h.authUseCase != nil {
		if accessToken != "" {
			if err := h.authUseCase.Logout(ctx, accessToken); err != nil && err != domain.ErrInvalidToken {
				h.respondAuthError(w, err, authOperationLogout)
				return
			}
		}

		if refreshToken != "" && refreshToken != accessToken {
			if err := h.authUseCase.Logout(ctx, refreshToken); err != nil && err != domain.ErrInvalidToken {
				h.respondAuthError(w, err, authOperationLogout)
				return
			}
		}
	}

	http.SetCookie(w, clearAuthCookie(r))
	http.SetCookie(w, clearRefreshCookie(r))

	RespondJSON(w, http.StatusOK, map[string]string{"message": "sessao encerrada"})
}

func (h *Handler) respondAuthError(w http.ResponseWriter, err error, operation authOperation) {
	switch err {
	case domain.ErrAuthUnavailable:
		RespondError(w, http.StatusServiceUnavailable, "auth service unavailable", "AUTH_UNAVAILABLE")
	case domain.ErrInvalidInput:
		RespondError(w, http.StatusBadRequest, "email e password obrigatorios", "INVALID_REQUEST")
	case domain.ErrUserAlreadyExists:
		RespondError(w, http.StatusConflict, "usuario ja existe", "ALREADY_EXISTS")
	case domain.ErrInvalidCredentials, domain.ErrUnauthorized:
		RespondError(w, http.StatusUnauthorized, "credenciais invalidas", "AUTH_ERROR")
	case domain.ErrInvalidToken:
		RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
	default:
		message := "falha na autenticacao"
		if operation == authOperationLogout {
			message = "falha ao encerrar sessao"
		}
		RespondError(w, http.StatusInternalServerError, message, "AUTH_ERROR")
	}
}

func decodeAuthJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthPayloadBytes)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
		return errors.New("payload invalido")
	}

	return nil
}

func buildAuthCookie(r *http.Request, token string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   requestUsesHTTPS(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func buildRefreshCookie(r *http.Request, token string) *http.Cookie {
	return &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		Secure:   requestUsesHTTPS(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func clearAuthCookie(r *http.Request) *http.Cookie {
	return &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   requestUsesHTTPS(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func clearRefreshCookie(r *http.Request) *http.Cookie {
	return &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   requestUsesHTTPS(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func extractRefreshTokenFromRequest(r *http.Request) (string, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader != "" {
		return parseBearerToken(authHeader)
	}

	if cookie, err := r.Cookie("refresh_token"); err == nil {
		if tokenString := strings.TrimSpace(cookie.Value); tokenString != "" {
			return tokenString, nil
		}
	}

	if cookie, err := r.Cookie("auth_token"); err == nil {
		if tokenString := strings.TrimSpace(cookie.Value); tokenString != "" {
			return tokenString, nil
		}
	}

	return "", nil
}

func requestUsesHTTPS(r *http.Request) bool {
	if r != nil && r.TLS != nil {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")), "https")
}
