package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/adapters/http/dto"
	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	"pokedex-platform/core/bff/mobile-bff/internal/ports"
)

type Handler struct {
	pokemonUseCase  ports.PokemonUseCase
	favoriteUseCase ports.FavoriteUseCase
	authUseCase     ports.AuthUseCase
	responseBuilder *ResponseBuilder
}

const maxAuthPayloadBytes int64 = 8 * 1024

func NewHandler(
	pokemonUseCase ports.PokemonUseCase,
	favoriteUseCase ports.FavoriteUseCase,
	authUseCase ports.AuthUseCase,
) *Handler {
	return &Handler{
		pokemonUseCase:  pokemonUseCase,
		favoriteUseCase: favoriteUseCase,
		authUseCase:     authUseCase,
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
	mux.HandleFunc("GET /api/v1/me", h.GetMe)
	mux.HandleFunc("GET /api/v1/me/favorites", h.GetUserFavorites)
	mux.HandleFunc("GET /api/v1/pokemons", h.ListPokemons)
	mux.HandleFunc("GET /api/v1/pokemons/search", h.SearchPokemons)
	mux.HandleFunc("GET /api/v1/pokemons/{id}/details", h.GetPokemonDetails)
	mux.HandleFunc("GET /api/v1/home", h.GetHome)
	mux.HandleFunc("GET /api/v1/regions", h.GetRegions)
	mux.HandleFunc("POST /api/v1/pokemons/{id}/favorite", h.RequireAuth(h.AddFavorite))
	mux.HandleFunc("DELETE /api/v1/pokemons/{id}/favorite", h.RequireAuth(h.RemoveFavorite))
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userID := getUserIDFromContext(ctx)
	response := h.responseBuilder.BuildProfileResponse(userID != "", getUserEmailFromContext(ctx))
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

type authOperation string

const (
	authOperationSignup  authOperation = "signup"
	authOperationLogin   authOperation = "login"
	authOperationRefresh authOperation = "refresh"
	authOperationLogout  authOperation = "logout"
)

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

	detail, err := h.pokemonUseCase.GetPokemonScreenDetails(ctx, pokemonID, userID)
	if err != nil {
		if err == domain.ErrPokemonNotFound {
			RespondError(w, http.StatusNotFound, "pokemon nao encontrado", "NOT_FOUND")
			return
		}
		RespondError(w, http.StatusInternalServerError, "falha ao obter detalhes do pokemon", "INTERNAL_ERROR")
		return
	}

	isFavorite := false
	if userID != "" {
		favoriteSet := h.buildFavoriteSet(ctx, userID)
		_, isFavorite = favoriteSet[normalizePokemonID(detail.Number)]
	}

	detailDTO := h.responseBuilder.BuildPokemonDetailScreenResponse(detail, isFavorite)
	RespondJSON(w, http.StatusOK, detailDTO)
}

func (h *Handler) GetHome(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	page, hasPage := getOptionalQueryParamInt(r, "page", 0)
	pageSize, hasSize := getOptionalQueryParamInt(r, "size", 20)
	userID := getUserIDFromContext(ctx)
	searchValue := strings.TrimSpace(r.URL.Query().Get("q"))
	selectedType := strings.TrimSpace(r.URL.Query().Get("type"))
	selectedOrdering := strings.TrimSpace(r.URL.Query().Get("order"))
	selectedRegion := strings.TrimSpace(r.URL.Query().Get("region"))
	paginate := hasPage

	if page < 0 {
		page = 0
	}

	if paginate {
		if !hasSize {
			pageSize = 20
		}
		if pageSize < 1 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}
	}

	pokemonPage, err := h.loadHomePokemonPage(ctx, userID, searchValue, selectedType)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao obter dados da home", "INTERNAL_ERROR")
		return
	}

	types, err := h.pokemonUseCase.ListTypes(ctx)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao obter filtros da home", "INTERNAL_ERROR")
		return
	}

	regions, err := h.pokemonUseCase.ListRegions(ctx)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao obter regioes da home", "INTERNAL_ERROR")
		return
	}

	filterHomePokemonByRegion(pokemonPage, selectedRegion)
	sortHomePokemonPage(pokemonPage, selectedOrdering)

	if paginate {
		paginateHomePokemonPage(pokemonPage, page, pageSize)
	}

	favoriteSet := h.buildFavoriteSet(ctx, userID)
	response := h.responseBuilder.BuildHomePageResponseWithTypes(
		pokemonPage,
		types,
		regions,
		favoriteSet,
		searchValue,
		selectedOrDefault(selectedType, "Todos os tipos"),
		selectedOrDefault(selectedOrdering, "Menor número"),
		selectedRegion,
	)
	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) GetRegions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	regions, err := h.pokemonUseCase.ListRegions(ctx)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao listar regioes", "INTERNAL_ERROR")
		return
	}

	RespondJSON(w, http.StatusOK, h.responseBuilder.BuildRegionsResponse(regions))
}

func (h *Handler) enrichFavoriteFlags(ctx context.Context, userID string, response *dto.RichPokemonListResponse) {
	if userID == "" || response == nil {
		return
	}

	favorites, err := h.favoriteUseCase.GetUserFavorites(ctx, userID)
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

func (h *Handler) buildFavoriteSet(ctx context.Context, userID string) map[string]struct{} {
	if userID == "" {
		return nil
	}

	favorites, err := h.favoriteUseCase.GetUserFavorites(ctx, userID)
	if err != nil {
		return nil
	}

	favoriteSet := make(map[string]struct{}, len(favorites))
	for _, id := range favorites {
		favoriteSet[normalizePokemonID(id)] = struct{}{}
	}

	return favoriteSet
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

	err := h.favoriteUseCase.AddFavorite(ctx, userID, pokemonID)
	if err == domain.ErrFavoriteAlreadyExists {
		RespondError(w, http.StatusConflict, "pokemon ja esta nos favoritos", "ALREADY_EXISTS")
		return
	}
	if err == domain.ErrPokemonNotFound {
		RespondError(w, http.StatusNotFound, "pokemon nao encontrado", "NOT_FOUND")
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

	err := h.favoriteUseCase.RemoveFavorite(ctx, userID, pokemonID)
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
		RespondJSON(w, http.StatusOK, h.responseBuilder.BuildFavoritesResponse(nil, nil, false))
		return
	}

	favorites, err := h.favoriteUseCase.GetUserFavorites(ctx, userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao listar favoritos", "INTERNAL_ERROR")
		return
	}

	items := make([]domain.Pokemon, 0, len(favorites))
	favoriteSet := make(map[string]struct{}, len(favorites))
	for _, favoriteID := range favorites {
		favoriteSet[normalizePokemonID(favoriteID)] = struct{}{}
		pokemon, err := h.pokemonUseCase.GetPokemonScreenDetails(ctx, favoriteID, userID)
		if err != nil {
			continue
		}
		items = append(items, domain.Pokemon{
			ID:           pokemon.ID,
			Name:         pokemon.Name,
			Number:       pokemon.Number,
			Types:        mapScreenTypesToNames(pokemon.Types),
			ImageURL:     pokemon.ImageURL,
			ElementColor: pokemon.ElementColor,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Number < items[j].Number
	})

	page := &domain.PokemonPage{Content: items}
	RespondJSON(w, http.StatusOK, h.responseBuilder.BuildFavoritesResponse(page, favoriteSet, true))
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

func getOptionalQueryParamInt(r *http.Request, key string, defaultVal int) (int, bool) {
	values, exists := r.URL.Query()[key]
	if !exists || len(values) == 0 {
		return defaultVal, false
	}

	trimmed := strings.TrimSpace(values[0])
	if trimmed == "" {
		return defaultVal, false
	}

	intVal, err := strconv.Atoi(trimmed)
	if err != nil {
		return defaultVal, true
	}

	return intVal, true
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

func (h *Handler) loadHomePokemonPage(
	ctx context.Context,
	userID string,
	searchValue string,
	selectedType string,
) (*domain.PokemonPage, error) {
	const fetchSize = 100
	items := make([]domain.Pokemon, 0, fetchSize)

	for page := 0; ; page++ {
		pokemonPage, err := h.pokemonUseCase.ListPokemons(ctx, page, fetchSize, userID)
		if err != nil {
			return nil, err
		}
		if pokemonPage == nil || len(pokemonPage.Content) == 0 {
			break
		}

		items = append(items, pokemonPage.Content...)
		if !pokemonPage.HasNext {
			break
		}
	}

	filtered := make([]domain.Pokemon, 0, len(items))
	searchTerm := strings.ToLower(strings.TrimSpace(searchValue))
	selectedType = strings.TrimSpace(selectedType)

	for _, pokemon := range items {
		if searchTerm != "" {
			name := strings.ToLower(strings.TrimSpace(pokemon.Name))
			number := normalizePokemonID(pokemon.Number)
			if !strings.Contains(name, searchTerm) && !strings.Contains(number, searchTerm) {
				continue
			}
		}

		if selectedType != "" && selectedType != "Todos os tipos" && !hasPokemonType(pokemon.Types, selectedType) {
			continue
		}

		filtered = append(filtered, pokemon)
	}

	return &domain.PokemonPage{Content: filtered}, nil
}

func paginateHomePokemonPage(page *domain.PokemonPage, currentPage, pageSize int) {
	if page == nil {
		return
	}

	if currentPage < 0 {
		currentPage = 0
	}
	if pageSize < 1 {
		pageSize = 20
	}

	start := currentPage * pageSize
	if start >= len(page.Content) {
		page.Content = []domain.Pokemon{}
		return
	}

	end := start + pageSize
	if end > len(page.Content) {
		end = len(page.Content)
	}

	page.Content = page.Content[start:end]
}

func sortHomePokemonPage(page *domain.PokemonPage, selectedOrdering string) {
	if page == nil {
		return
	}

	switch selectedOrdering {
	case "Maior número":
		sort.Slice(page.Content, func(i, j int) bool { return page.Content[i].Number > page.Content[j].Number })
	case "A-Z":
		sort.Slice(page.Content, func(i, j int) bool { return page.Content[i].Name < page.Content[j].Name })
	case "Z-A":
		sort.Slice(page.Content, func(i, j int) bool { return page.Content[i].Name > page.Content[j].Name })
	default:
		sort.Slice(page.Content, func(i, j int) bool { return page.Content[i].Number < page.Content[j].Number })
	}
}

func filterHomePokemonByRegion(page *domain.PokemonPage, region string) {
	if page == nil || strings.TrimSpace(region) == "" {
		return
	}

	filtered := make([]domain.Pokemon, 0, len(page.Content))
	for _, pokemon := range page.Content {
		if matchesRegion(pokemon.Number, region) {
			filtered = append(filtered, pokemon)
		}
	}
	page.Content = filtered
}

func matchesRegion(number string, region string) bool {
	parsed, ok := parsePokemonNumber(number)
	if !ok {
		return true
	}

	switch strings.ToLower(strings.TrimSpace(region)) {
	case "kanto":
		return parsed >= 1 && parsed <= 151
	case "johto":
		return parsed >= 152 && parsed <= 251
	case "hoenn":
		return parsed >= 252 && parsed <= 386
	case "sinnoh":
		return parsed >= 387 && parsed <= 493
	case "unova":
		return parsed >= 494 && parsed <= 649
	case "kalos":
		return parsed >= 650 && parsed <= 721
	case "alola":
		return parsed >= 722 && parsed <= 809
	case "galar":
		return parsed >= 810 && parsed <= 905
	default:
		return true
	}
}

func parsePokemonNumber(number string) (int, bool) {
	normalized := normalizePokemonID(number)
	parsed, err := strconv.Atoi(normalized)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func selectedOrDefault(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func mapScreenTypesToNames(types []domain.Type) []string {
	result := make([]string, len(types))
	for i, item := range types {
		result[i] = item.Name
	}
	return result
}

func hasPokemonType(types []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, item := range types {
		if strings.EqualFold(strings.TrimSpace(item), target) {
			return true
		}
	}
	return false
}
