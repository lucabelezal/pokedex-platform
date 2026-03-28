package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

func NewMux(authService *service.AuthService) *http.ServeMux {
	h := &Handler{authService: authService}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("POST /v1/auth/signup", h.signup)
	mux.HandleFunc("POST /v1/auth/login", h.login)
	mux.HandleFunc("POST /v1/auth/refresh", h.refresh)
	mux.HandleFunc("POST /v1/auth/logout", h.logout)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, healthResponse{Status: "ok", Service: "auth-service"})
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "payload invalido")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := h.authService.Signup(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			respondError(w, http.StatusBadRequest, "email ou senha invalidos")
		case errors.Is(err, repository.ErrUserAlreadyExist):
			respondError(w, http.StatusConflict, "usuario ja existe")
		default:
			respondError(w, http.StatusInternalServerError, "falha ao criar usuario")
		}
		return
	}

	respondJSON(w, http.StatusCreated, result)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "payload invalido")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			respondError(w, http.StatusBadRequest, "email ou senha invalidos")
		case errors.Is(err, service.ErrInvalidCredentials):
			respondError(w, http.StatusUnauthorized, "credenciais invalidas")
		default:
			respondError(w, http.StatusInternalServerError, "falha ao autenticar")
		}
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := extractBearerToken(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "token invalido")
		return
	}

	result, err := h.authService.Refresh(tokenString)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "token invalido")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	tokenString, err := extractBearerToken(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "token invalido")
		return
	}

	if err := h.authService.Logout(tokenString); err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao encerrar sessao")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "sessao encerrada"})
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
