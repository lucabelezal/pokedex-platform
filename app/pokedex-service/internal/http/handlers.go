package http

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type pingResponse struct {
	Message string `json:"message"`
}

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/v1/pokemon/ping", pingHandler)
	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, healthResponse{Status: "ok", Service: "pokedex-service"})
}

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, pingResponse{Message: "service is alive"})
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
