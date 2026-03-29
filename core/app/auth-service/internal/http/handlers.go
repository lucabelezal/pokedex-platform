package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"pokedex-platform/core/app/auth-service/internal/repository"
	"pokedex-platform/core/app/auth-service/internal/service"
)

type Handler struct {
	authService *service.AuthService
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type introspectResponse struct {
	Active bool `json:"active"`
}

const maxAuthRequestBodyBytes int64 = 8 * 1024

func NewMux(authService *service.AuthService) *http.ServeMux {
	h := &Handler{authService: authService}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("POST /v1/auth/signup", h.signup)
	mux.HandleFunc("POST /v1/auth/login", h.login)
	mux.HandleFunc("POST /v1/auth/refresh", h.refresh)
	mux.HandleFunc("POST /v1/auth/logout", h.logout)
	mux.HandleFunc("POST /v1/auth/introspect", h.introspect)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, healthResponse{Status: "ok", Service: "auth-service"})
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		logAuthAudit(r, "signup", "invalid_request", "reason=payload_invalido")
		respondError(w, http.StatusBadRequest, "payload invalido")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := h.authService.Signup(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			logAuthAudit(r, "signup", "rejected", "email=%q", strings.ToLower(strings.TrimSpace(req.Email)))
			respondError(w, http.StatusBadRequest, "email invalido ou senha com menos de 8 caracteres")
		case errors.Is(err, repository.ErrUserAlreadyExist):
			logAuthAudit(r, "signup", "conflict", "email=%q", strings.ToLower(strings.TrimSpace(req.Email)))
			respondError(w, http.StatusConflict, "usuario ja existe")
		default:
			logAuthAudit(r, "signup", "error", "email=%q", strings.ToLower(strings.TrimSpace(req.Email)))
			respondError(w, http.StatusInternalServerError, "falha ao criar usuario")
		}
		return
	}

	logAuthAudit(r, "signup", "success", "user_id=%q email=%q", result.UserID, result.Email)

	respondJSON(w, http.StatusCreated, result)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		logAuthAudit(r, "login", "invalid_request", "reason=payload_invalido")
		respondError(w, http.StatusBadRequest, "payload invalido")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			logAuthAudit(r, "login", "rejected", "email=%q", strings.ToLower(strings.TrimSpace(req.Email)))
			respondError(w, http.StatusBadRequest, "email invalido ou senha obrigatoria")
		case errors.Is(err, service.ErrInvalidCredentials):
			logAuthAudit(r, "login", "invalid_credentials", "email=%q", strings.ToLower(strings.TrimSpace(req.Email)))
			respondError(w, http.StatusUnauthorized, "credenciais invalidas")
		default:
			logAuthAudit(r, "login", "error", "email=%q", strings.ToLower(strings.TrimSpace(req.Email)))
			respondError(w, http.StatusInternalServerError, "falha ao autenticar")
		}
		return
	}

	logAuthAudit(r, "login", "success", "user_id=%q email=%q", result.UserID, result.Email)

	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := extractBearerToken(r)
	if err != nil {
		logAuthAudit(r, "refresh", "invalid_request", "reason=token_ausente")
		respondError(w, http.StatusUnauthorized, "token invalido")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := h.authService.Refresh(ctx, tokenString)
	if err != nil {
		logAuthAudit(r, "refresh", "invalid_token")
		respondError(w, http.StatusUnauthorized, "token invalido")
		return
	}

	logAuthAudit(r, "refresh", "success", "user_id=%q email=%q", result.UserID, result.Email)

	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	tokenString, err := extractBearerToken(r)
	if err != nil {
		logAuthAudit(r, "logout", "invalid_request", "reason=token_ausente")
		respondError(w, http.StatusUnauthorized, "token invalido")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.authService.Logout(ctx, tokenString); err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			logAuthAudit(r, "logout", "invalid_token")
			respondError(w, http.StatusUnauthorized, "token invalido")
			return
		}

		logAuthAudit(r, "logout", "error")
		respondError(w, http.StatusInternalServerError, "falha ao encerrar sessao")
		return
	}

	logAuthAudit(r, "logout", "success")

	respondJSON(w, http.StatusOK, map[string]string{"message": "sessao encerrada"})
}

func (h *Handler) introspect(w http.ResponseWriter, r *http.Request) {
	tokenString, err := extractBearerToken(r)
	if err != nil {
		respondJSON(w, http.StatusOK, introspectResponse{Active: false})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	active, err := h.authService.IsAccessTokenActive(ctx, tokenString)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao validar token")
		return
	}

	respondJSON(w, http.StatusOK, introspectResponse{Active: active})
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization ausente")
	}

	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", errors.New("authorization invalido")
	}

	tokenString := authHeader[len(prefix):]
	if tokenString == "" {
		return "", errors.New("token vazio")
	}

	return tokenString, nil
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthRequestBodyBytes)
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

func logAuthAudit(r *http.Request, action, outcome string, args ...any) {
	clientIP := authClientIP(r)
	message := ""
	if len(args) > 0 {
		format, _ := args[0].(string)
		message = fmt.Sprintf(format, args[1:]...)
	}

	if message != "" {
		log.Printf("auth_audit action=%s outcome=%s client_ip=%q %s", action, outcome, clientIP, message)
		return
	}

	log.Printf("auth_audit action=%s outcome=%s client_ip=%q", action, outcome, clientIP)
}

func authClientIP(r *http.Request) string {
	forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwardedFor != "" {
		first, _, _ := strings.Cut(forwardedFor, ",")
		if candidate := strings.TrimSpace(first); candidate != "" {
			return candidate
		}
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
