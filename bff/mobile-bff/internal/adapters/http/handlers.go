package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pokedex-platform/bff/mobile-bff/internal/adapters/http/dto"
	"pokedex-platform/bff/mobile-bff/internal/adapters/repository"
	"pokedex-platform/bff/mobile-bff/internal/domain"
	"pokedex-platform/bff/mobile-bff/internal/ports"
)

type Handler struct {
	pokemonUseCase  ports.PokemonUseCase
	favoriteUseCase ports.FavoriteUseCase
	authClient      *repository.AuthServiceClient
	favoriteRepo    ports.FavoriteRepository
	responseBuilder *ResponseBuilder
}

func NewHandler(
	pokemonUseCase ports.PokemonUseCase,
	favoriteUseCase ports.FavoriteUseCase,
	authClient *repository.AuthServiceClient,
	favoriteRepo ports.FavoriteRepository,
) *Handler {
	return &Handler{
		pokemonUseCase:  pokemonUseCase,
		favoriteUseCase: favoriteUseCase,
		authClient:      authClient,
		favoriteRepo:    favoriteRepo,
		responseBuilder: NewResponseBuilder(),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/v1/health", h.Health)
	mux.HandleFunc("POST /api/v1/auth/signup", h.Signup)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", h.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", h.Logout)
	mux.HandleFunc("GET /api/v1/me", h.withAuth(h.GetMe))
	mux.HandleFunc("GET /api/v1/me/favorites", h.withAuth(h.GetUserFavorites))
	mux.HandleFunc("GET /api/v1/pokemons", h.ListPokemons)
	mux.HandleFunc("GET /api/v1/pokemons/search", h.SearchPokemons)
	mux.HandleFunc("GET /api/v1/pokemons/{id}/details", h.GetPokemonDetails)
	mux.HandleFunc("GET /api/v1/home", h.GetHome)
	mux.HandleFunc("POST /api/v1/pokemons/{id}/favorite", h.RequireAuth(h.AddFavorite))
	mux.HandleFunc("DELETE /api/v1/pokemons/{id}/favorite", h.RequireAuth(h.RemoveFavorite))
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userID := getUserIDFromContext(ctx)
	if userID == "" {
		RespondError(w, http.StatusUnauthorized, "autenticacao obrigatoria", "UNAUTHORIZED")
		return
	}

	response := struct {
		Authenticated bool   `json:"authenticated"`
		UserID        string `json:"user_id"`
		Email         string `json:"email,omitempty"`
	}{
		Authenticated: true,
		UserID:        userID,
		Email:         getUserEmailFromContext(ctx),
	}

	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	health := h.responseBuilder.BuildHealthResponse()
	RespondJSON(w, http.StatusOK, health)
}

// Signup gerencia registro de usuário
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "email e password obrigatorios", "INVALID_REQUEST")
		return
	}

	if req.Email == "" || req.Password == "" {
		RespondError(w, http.StatusBadRequest, "email e password obrigatorios", "INVALID_REQUEST")
		return
	}

	authResp, err := h.authClient.Signup(ctx, req.Email, req.Password)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error(), "AUTH_ERROR")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authResp.AccessToken,
		Path:     "/",
		MaxAge:   authResp.ExpiresIn,
		Secure:   false, // definir como true em produção com HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	response := struct {
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
	}

	RespondJSON(w, http.StatusCreated, response)
}

// Login gerencia autenticação de usuário
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "email e password obrigatorios", "INVALID_REQUEST")
		return
	}

	if req.Email == "" || req.Password == "" {
		RespondError(w, http.StatusBadRequest, "email e password obrigatorios", "INVALID_REQUEST")
		return
	}

	authResp, err := h.authClient.Login(ctx, req.Email, req.Password)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err.Error(), "AUTH_ERROR")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authResp.AccessToken,
		Path:     "/",
		MaxAge:   authResp.ExpiresIn,
		Secure:   false, // definir como true em produção com HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	response := struct {
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
	}

	RespondJSON(w, http.StatusOK, response)
}

// Refresh renova o token de acesso do usuário autenticado
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "metodo nao permitido", "METHOD_NOT_ALLOWED")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tokenString, err := extractTokenFromRequest(r)
	if err != nil || tokenString == "" {
		RespondError(w, http.StatusUnauthorized, "token invalido", "INVALID_TOKEN")
		return
	}

	authResp, err := h.authClient.Refresh(ctx, tokenString)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err.Error(), "AUTH_ERROR")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authResp.AccessToken,
		Path:     "/",
		MaxAge:   authResp.ExpiresIn,
		Secure:   false, // definir como true em produção com HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	response := struct {
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
	}

	RespondJSON(w, http.StatusOK, response)
}

// Logout encerra a sessão e remove o cookie de autenticação
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "metodo nao permitido", "METHOD_NOT_ALLOWED")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tokenString, _ := extractTokenFromRequest(r)
	if tokenString != "" {
		_ = h.authClient.Logout(ctx, tokenString)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false, // definir como true em produção com HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	RespondJSON(w, http.StatusOK, map[string]string{"message": "sessao encerrada"})
}

func (h *Handler) ListPokemons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	page := getQueryParamInt(r, "page", 0)
	pageSize := getQueryParamInt(r, "size", 20)
	userID := getUserIDFromContext(ctx)
	typeFilter := r.URL.Query().Get("type")

	var (
		pokemonPage *domain.PokemonPage
		err         error
	)

	if typeFilter != "" {
		pokemonPage, err = h.pokemonUseCase.FilterByType(ctx, typeFilter, page, pageSize, userID)
	} else {
		pokemonPage, err = h.pokemonUseCase.ListPokemons(ctx, page, pageSize, userID)
	}

	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao listar pokemons", "INTERNAL_ERROR")
		return
	}

	response := h.responseBuilder.BuildRichPokemonListResponse(pokemonPage)
	h.enrichFavoriteFlags(ctx, userID, response)
	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) SearchPokemons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query().Get("q")
	if query == "" {
		RespondError(w, http.StatusBadRequest, "termo de busca obrigatorio", "INVALID_REQUEST")
		return
	}

	page := getQueryParamInt(r, "page", 0)
	pageSize := getQueryParamInt(r, "size", 20)
	userID := getUserIDFromContext(ctx)

	pokemonPage, err := h.pokemonUseCase.SearchPokemons(ctx, query, page, pageSize, userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao buscar pokemons", "INTERNAL_ERROR")
		return
	}

	response := h.responseBuilder.BuildRichPokemonListResponse(pokemonPage)
	h.enrichFavoriteFlags(ctx, userID, response)
	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) GetPokemonDetails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pokemonID := r.PathValue("id")
	if pokemonID == "" {
		RespondError(w, http.StatusBadRequest, "id do pokemon obrigatorio", "INVALID_REQUEST")
		return
	}

	userID := getUserIDFromContext(ctx)

	detail, err := h.pokemonUseCase.GetPokemonDetails(ctx, pokemonID, userID)
	if err != nil {
		if err == domain.ErrPokemonNotFound {
			RespondError(w, http.StatusNotFound, "pokemon nao encontrado", "NOT_FOUND")
			return
		}
		RespondError(w, http.StatusInternalServerError, "falha ao obter detalhes do pokemon", "INTERNAL_ERROR")
		return
	}

	detailDTO := h.responseBuilder.BuildPokemonDetailDTO(detail)
	RespondJSON(w, http.StatusOK, detailDTO)
}

func (h *Handler) GetHome(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	page := getQueryParamInt(r, "page", 0)
	pageSize := getQueryParamInt(r, "size", 20)
	userID := getUserIDFromContext(ctx)

	pokemonPage, err := h.pokemonUseCase.GetHomeData(ctx, page, pageSize, userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao obter dados da home", "INTERNAL_ERROR")
		return
	}

	response := h.responseBuilder.BuildHomePageResponse(pokemonPage)
	h.enrichFavoriteFlags(ctx, userID, response.Data)
	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) enrichFavoriteFlags(ctx context.Context, userID string, response *dto.RichPokemonListResponse) {
	if userID == "" || response == nil || h.favoriteRepo == nil {
		return
	}

	favorites, err := h.favoriteRepo.GetUserFavorites(ctx, userID)
	if err != nil {
		return
	}

	favoriteSet := make(map[string]struct{}, len(favorites))
	for _, id := range favorites {
		favoriteSet[normalizePokemonID(id)] = struct{}{}
	}

	for i := range response.Content {
		_, isFavorite := favoriteSet[normalizePokemonID(response.Content[i].Number)]
		response.Content[i].IsFavorite = isFavorite
	}
}

func normalizePokemonID(value string) string {
	normalized := strings.TrimLeft(strings.TrimSpace(value), "0")
	if normalized == "" {
		return "0"
	}
	return normalized
}

func (h *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pokemonID := r.PathValue("id")
	if pokemonID == "" {
		RespondError(w, http.StatusBadRequest, "id do pokemon obrigatorio", "INVALID_REQUEST")
		return
	}

	userID := getUserIDFromContext(ctx)

	err := h.favoriteRepo.AddFavorite(ctx, userID, pokemonID)
	if err == domain.ErrFavoriteAlreadyExists {
		RespondError(w, http.StatusConflict, "pokemon ja esta nos favoritos", "ALREADY_EXISTS")
		return
	}
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao adicionar favorito", "INTERNAL_ERROR")
		return
	}

	response := dto.FavoriteResponse{
		Message:    "Pokemon adicionado aos favoritos",
		PokemonID:  pokemonID,
		IsFavorite: true,
	}
	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pokemonID := r.PathValue("id")
	if pokemonID == "" {
		RespondError(w, http.StatusBadRequest, "id do pokemon obrigatorio", "INVALID_REQUEST")
		return
	}

	userID := getUserIDFromContext(ctx)

	err := h.favoriteRepo.RemoveFavorite(ctx, userID, pokemonID)
	if err == domain.ErrFavoriteNotFound {
		RespondError(w, http.StatusNotFound, "favorito nao encontrado", "NOT_FOUND")
		return
	}
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao remover favorito", "INTERNAL_ERROR")
		return
	}

	response := dto.MessageResponse{
		Message: "Pokemon removido dos favoritos",
	}
	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userID := getUserIDFromContext(ctx)
	if userID == "" {
		RespondError(w, http.StatusUnauthorized, "autenticacao obrigatoria", "UNAUTHORIZED")
		return
	}

	favorites, err := h.favoriteRepo.GetUserFavorites(ctx, userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao listar favoritos", "INTERNAL_ERROR")
		return
	}

	response := struct {
		Favorites []string `json:"favorites"`
		Count     int      `json:"count"`
	}{
		Favorites: favorites,
		Count:     len(favorites),
	}
	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) withAuth(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())
		if userID == "" {
			RespondError(w, http.StatusUnauthorized, "autenticacao obrigatoria", "UNAUTHORIZED")
			return
		}
		handler(w, r)
	}
}

func getQueryParamInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return intVal
}

// RequireAuth envolve um handler para exigir autenticação
func (h *Handler) RequireAuth(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())
		if userID == "" {
			RespondError(w, http.StatusUnauthorized, "autenticacao obrigatoria", "UNAUTHORIZED")
			return
		}
		handler(w, r)
	}
}
