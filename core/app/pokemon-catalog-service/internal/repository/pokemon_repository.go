package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"pokedex-platform/core/app/pokemon-catalog-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPokemonNotFound = errors.New("pokemon nao encontrado")
var ErrFavoriteAlreadyExists = errors.New("favorito ja existe")
var ErrFavoriteNotFound = errors.New("favorito nao encontrado")

type PokemonRepository interface {
	GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error)
	Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error)
	GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error)
	GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
	GetByIDs(ctx context.Context, ids []string) ([]domain.Pokemon, error)
	GetDetailByID(ctx context.Context, id string) (*domain.PokemonDetail, error)
	ListTypes(ctx context.Context) ([]domain.Type, error)
	ListRegions(ctx context.Context) ([]domain.Region, error)
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
}

// DBPool abstrai o pool pgx para permitir mocks em testes.
type DBPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
}

type PostgresPokemonRepository struct {
	db DBPool
}

func NewPostgresPokemonRepository(db *pgxpool.Pool) *PostgresPokemonRepository {
	return &PostgresPokemonRepository{db: db}
}

// NewPostgresPokemonRepositoryWithPool cria repositório aceitando qualquer DBPool (incluindo mocks).
func NewPostgresPokemonRepositoryWithPool(db DBPool) *PostgresPokemonRepository {
	return &PostgresPokemonRepository{db: db}
}

func (r *PostgresPokemonRepository) GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error) {
	page, pageSize = sanitizePage(page, pageSize)
	offset := page * pageSize

	query := `
		SELECT
			p.id::text,
			p.name,
			p.number,
			COALESCE(array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL), ARRAY[]::text[]),
			COALESCE(p.height::double precision, 0),
			COALESCE(p.weight::double precision, 0),
			COALESCE(p.description, ''),
			COALESCE(
				jsonb_extract_path_text(p.sprites, 'other', 'home', 'front_default'),
				jsonb_extract_path_text(p.sprites, 'other', 'official-artwork', 'front_default'),
				p.sprites->>'front_default',
				''
			),
			COALESCE((array_agg(t.color ORDER BY t.id) FILTER (WHERE t.color IS NOT NULL))[1], '#A9AC86'),
			COALESCE((array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL))[1], ''),
			now(),
			now()
		FROM pokemons p
		LEFT JOIN pokemon_types pt ON pt.pokemon_id = p.id
		LEFT JOIN types t ON t.id = pt.type_id
		GROUP BY p.id, p.name, p.number, p.height, p.weight, p.description, p.sprites
		ORDER BY p.number ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := readRows(rows)
	if err != nil {
		return nil, err
	}

	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pokemons").Scan(&total); err != nil {
		return nil, err
	}

	return buildPage(items, total, page, pageSize), nil
}

func (r *PostgresPokemonRepository) Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error) {
	page, pageSize = sanitizePage(page, pageSize)
	offset := page * pageSize
	q := "%" + query + "%"

	sql := `
		SELECT
			p.id::text,
			p.name,
			p.number,
			COALESCE(array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL), ARRAY[]::text[]),
			COALESCE(p.height::double precision, 0),
			COALESCE(p.weight::double precision, 0),
			COALESCE(p.description, ''),
			COALESCE(
				jsonb_extract_path_text(p.sprites, 'other', 'home', 'front_default'),
				jsonb_extract_path_text(p.sprites, 'other', 'official-artwork', 'front_default'),
				p.sprites->>'front_default',
				''
			),
			COALESCE((array_agg(t.color ORDER BY t.id) FILTER (WHERE t.color IS NOT NULL))[1], '#A9AC86'),
			COALESCE((array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL))[1], ''),
			now(),
			now()
		FROM pokemons p
		LEFT JOIN pokemon_types pt ON pt.pokemon_id = p.id
		LEFT JOIN types t ON t.id = pt.type_id
		WHERE p.name ILIKE $1 OR p.number ILIKE $1
		GROUP BY p.id, p.name, p.number, p.height, p.weight, p.description, p.sprites
		ORDER BY p.number ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, sql, q, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := readRows(rows)
	if err != nil {
		return nil, err
	}

	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pokemons WHERE name ILIKE $1 OR number ILIKE $1", q).Scan(&total); err != nil {
		return nil, err
	}

	return buildPage(items, total, page, pageSize), nil
}

func (r *PostgresPokemonRepository) GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error) {
	page, pageSize = sanitizePage(page, pageSize)
	offset := page * pageSize

	query := `
		SELECT
			p.id::text,
			p.name,
			p.number,
			COALESCE(array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL), ARRAY[]::text[]),
			COALESCE(p.height::double precision, 0),
			COALESCE(p.weight::double precision, 0),
			COALESCE(p.description, ''),
			COALESCE(
				jsonb_extract_path_text(p.sprites, 'other', 'home', 'front_default'),
				jsonb_extract_path_text(p.sprites, 'other', 'official-artwork', 'front_default'),
				p.sprites->>'front_default',
				''
			),
			COALESCE((array_agg(t.color ORDER BY t.id) FILTER (WHERE t.color IS NOT NULL))[1], '#A9AC86'),
			COALESCE((array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL))[1], ''),
			now(),
			now()
		FROM pokemons p
		JOIN pokemon_types pt_filter ON pt_filter.pokemon_id = p.id
		JOIN types t_filter ON t_filter.id = pt_filter.type_id
		LEFT JOIN pokemon_types pt ON pt.pokemon_id = p.id
		LEFT JOIN types t ON t.id = pt.type_id
		WHERE t_filter.name = $1
		GROUP BY p.id, p.name, p.number, p.height, p.weight, p.description, p.sprites
		ORDER BY p.number ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, typeFilter, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := readRows(rows)
	if err != nil {
		return nil, err
	}

	var total int64
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT p.id)
		FROM pokemons p
		JOIN pokemon_types pt ON pt.pokemon_id = p.id
		JOIN types t ON t.id = pt.type_id
		WHERE t.name = $1
	`, typeFilter).Scan(&total); err != nil {
		return nil, err
	}

	return buildPage(items, total, page, pageSize), nil
}

func (r *PostgresPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
	query := `
		SELECT
			p.id::text,
			p.name,
			p.number,
			COALESCE(array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL), ARRAY[]::text[]),
			COALESCE(p.height::double precision, 0),
			COALESCE(p.weight::double precision, 0),
			COALESCE(p.description, ''),
			COALESCE(
				jsonb_extract_path_text(p.sprites, 'other', 'home', 'front_default'),
				jsonb_extract_path_text(p.sprites, 'other', 'official-artwork', 'front_default'),
				p.sprites->>'front_default',
				''
			),
			COALESCE((array_agg(t.color ORDER BY t.id) FILTER (WHERE t.color IS NOT NULL))[1], '#A9AC86'),
			COALESCE((array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL))[1], ''),
			now(),
			now()
		FROM pokemons p
		LEFT JOIN pokemon_types pt ON pt.pokemon_id = p.id
		LEFT JOIN types t ON t.id = pt.type_id
		WHERE p.id::text = $1 OR p.number = $1
		GROUP BY p.id, p.name, p.number, p.height, p.weight, p.description, p.sprites
	`

	var p domain.Pokemon
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Number,
		&p.Types,
		&p.Height,
		&p.Weight,
		&p.Description,
		&p.ImageURL,
		&p.ElementColor,
		&p.ElementType,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, ErrPokemonNotFound
	}

	return &p, nil
}

func (r *PostgresPokemonRepository) GetByIDs(ctx context.Context, ids []string) ([]domain.Pokemon, error) {
	if len(ids) == 0 {
		return []domain.Pokemon{}, nil
	}

	query := `
		SELECT
			p.id::text,
			p.name,
			p.number,
			COALESCE(array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL), ARRAY[]::text[]),
			COALESCE(p.height::double precision, 0),
			COALESCE(p.weight::double precision, 0),
			COALESCE(p.description, ''),
			COALESCE(
				jsonb_extract_path_text(p.sprites, 'other', 'home', 'front_default'),
				jsonb_extract_path_text(p.sprites, 'other', 'official-artwork', 'front_default'),
				p.sprites->>'front_default',
				''
			),
			COALESCE((array_agg(t.color ORDER BY t.id) FILTER (WHERE t.color IS NOT NULL))[1], '#A9AC86'),
			COALESCE((array_agg(t.name ORDER BY t.id) FILTER (WHERE t.name IS NOT NULL))[1], ''),
			now(),
			now()
		FROM pokemons p
		LEFT JOIN pokemon_types pt ON pt.pokemon_id = p.id
		LEFT JOIN types t ON t.id = pt.type_id
		WHERE p.id::text = ANY($1)
		GROUP BY p.id, p.name, p.number, p.height, p.weight, p.description, p.sprites
		ORDER BY p.id
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pokemons por ids: %w", err)
	}
	defer rows.Close()

	result := make([]domain.Pokemon, 0)
	for rows.Next() {
		var p domain.Pokemon
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Number,
			&p.Types,
			&p.Height,
			&p.Weight,
			&p.Description,
			&p.ImageURL,
			&p.ElementColor,
			&p.ElementType,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear pokemon: %w", err)
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar resultados: %w", err)
	}

	return result, nil
}

func (r *PostgresPokemonRepository) ListTypes(ctx context.Context) ([]domain.Type, error) {
	rows, err := r.db.Query(ctx, `
		SELECT name, color
		FROM types
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	types := make([]domain.Type, 0)
	for rows.Next() {
		var item domain.Type
		if err := rows.Scan(&item.Name, &item.Color); err != nil {
			return nil, err
		}
		types = append(types, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return types, nil
}

func (r *PostgresPokemonRepository) ListRegions(ctx context.Context) ([]domain.Region, error) {
	rows, err := r.db.Query(ctx, `
		SELECT r.id, r.name, g.id
		FROM regions r
		JOIN generations g ON g.region_id = r.id
		WHERE g.id <= 8
		ORDER BY g.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	regions := make([]domain.Region, 0)
	for rows.Next() {
		var (
			id           int64
			name         string
			generationID int64
		)
		if err := rows.Scan(&id, &name, &generationID); err != nil {
			return nil, err
		}

		regions = append(regions, domain.Region{
			ID:         strings.ToLower(strings.TrimSpace(name)),
			Name:       name,
			Generation: generationLabel(int(generationID)),
		})
	}

	return regions, rows.Err()
}

func (r *PostgresPokemonRepository) GetDetailByID(ctx context.Context, id string) (*domain.PokemonDetail, error) {
	base, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			COALESCE(NULLIF(s.species_pt, ''), NULLIF(s.species_en, ''), ''),
			COALESCE(p.gender_male, -1),
			COALESCE(p.gender_female, -1),
			COALESCE(r.name, ''),
			COALESCE(g.name, ''),
			COALESCE(p.evolution_chain_id, 0)
		FROM pokemons p
		LEFT JOIN species s ON s.id = p.species_id
		LEFT JOIN regions r ON r.id = p.region_id
		LEFT JOIN generations g ON g.id = p.generation_id
		WHERE p.id::text = $1 OR p.number = $1
	`

	var (
		category         string
		genderMale       float64
		genderFemale     float64
		region           string
		generation       string
		evolutionChainID int64
	)

	if err := r.db.QueryRow(ctx, query, id).Scan(
		&category,
		&genderMale,
		&genderFemale,
		&region,
		&generation,
		&evolutionChainID,
	); err != nil {
		return nil, err
	}

	abilities, err := r.listAbilityNames(ctx, id)
	if err != nil {
		return nil, err
	}

	weaknesses, err := r.listWeaknesses(ctx, id)
	if err != nil {
		return nil, err
	}

	evolutions, err := r.listEvolutions(ctx, base.Number, evolutionChainID)
	if err != nil {
		return nil, err
	}

	detail := &domain.PokemonDetail{
		ID:           base.ID,
		Name:         base.Name,
		Number:       base.Number,
		Types:        mapTypes(base.Types),
		Description:  normalizeDescription(base.Number, base.Description),
		ImageURL:     base.ImageURL,
		ElementColor: base.ElementColor,
		Height:       base.Height,
		Weight:       base.Weight,
		Category:     normalizeCategory(base.Number, category),
		Abilities:    normalizeAbilities(base.Number, abilities),
		Weaknesses:   weaknesses,
		Evolutions:   evolutions,
		Region:       region,
		Generation:   generation,
	}

	if genderMale >= 0 {
		detail.GenderMale = &genderMale
	}
	if genderFemale >= 0 {
		detail.GenderFemale = &genderFemale
	}

	return detail, nil
}

type InMemoryPokemonRepository struct {
	items []domain.Pokemon
}

func NewInMemoryPokemonRepository() *InMemoryPokemonRepository {
	return &InMemoryPokemonRepository{items: []domain.Pokemon{
		{
			ID:           "00000000-0000-0000-0000-000000000001",
			Name:         "Bulbasaur",
			Number:       "001",
			Types:        []string{"Grass", "Poison"},
			Height:       0.7,
			Weight:       6.9,
			Description:  "A seed Pokemon.",
			ImageURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png",
			ElementColor: "#78C850",
			ElementType:  "Grass",
		},
		{
			ID:           "00000000-0000-0000-0000-000000000004",
			Name:         "Charmander",
			Number:       "004",
			Types:        []string{"Fire"},
			Height:       0.6,
			Weight:       8.5,
			Description:  "A fire lizard Pokemon.",
			ImageURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/4.png",
			ElementColor: "#F08030",
			ElementType:  "Fire",
		},
	}}
}

func (r *InMemoryPokemonRepository) GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error) {
	_ = ctx
	page, pageSize = sanitizePage(page, pageSize)
	return pageSlice(r.items, page, pageSize), nil
}

func (r *InMemoryPokemonRepository) Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error) {
	_ = ctx
	q := strings.ToLower(strings.TrimSpace(query))
	filtered := make([]domain.Pokemon, 0, len(r.items))
	for _, p := range r.items {
		if strings.Contains(strings.ToLower(p.Name), q) || strings.Contains(strings.ToLower(p.Number), q) {
			filtered = append(filtered, p)
		}
	}
	page, pageSize = sanitizePage(page, pageSize)
	return pageSlice(filtered, page, pageSize), nil
}

func (r *InMemoryPokemonRepository) GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error) {
	_ = ctx
	target := strings.ToLower(strings.TrimSpace(typeFilter))
	filtered := make([]domain.Pokemon, 0, len(r.items))
	for _, p := range r.items {
		for _, t := range p.Types {
			if strings.ToLower(t) == target {
				filtered = append(filtered, p)
				break
			}
		}
	}
	page, pageSize = sanitizePage(page, pageSize)
	return pageSlice(filtered, page, pageSize), nil
}

func (r *InMemoryPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
	_ = ctx
	for _, p := range r.items {
		if p.ID == id || p.Number == id {
			cp := p
			return &cp, nil
		}
	}
	return nil, ErrPokemonNotFound
}

func (r *InMemoryPokemonRepository) GetByIDs(ctx context.Context, ids []string) ([]domain.Pokemon, error) {
	_ = ctx
	result := make([]domain.Pokemon, 0)
	for _, id := range ids {
		for _, p := range r.items {
			if p.ID == id || p.Number == id {
				result = append(result, p)
				break
			}
		}
	}
	return result, nil
}

func (r *InMemoryPokemonRepository) AddFavorite(ctx context.Context, userID, pokemonID string) error {
	return nil
}

func (r *InMemoryPokemonRepository) RemoveFavorite(ctx context.Context, userID, pokemonID string) error {
	return nil
}

func (r *InMemoryPokemonRepository) ListTypes(ctx context.Context) ([]domain.Type, error) {
	_ = ctx
	return []domain.Type{
		{Name: "Grama", Color: "#7AC74C"},
		{Name: "Venenoso", Color: "#A33EA1"},
		{Name: "Fogo", Color: "#EE8130"},
	}, nil
}

func (r *InMemoryPokemonRepository) ListRegions(ctx context.Context) ([]domain.Region, error) {
	_ = ctx
	return []domain.Region{
		{ID: "kanto", Name: "Kanto", Generation: "1º Geração"},
		{ID: "johto", Name: "Johto", Generation: "2º Geração"},
	}, nil
}

func (r *InMemoryPokemonRepository) GetDetailByID(ctx context.Context, id string) (*domain.PokemonDetail, error) {
	base, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	abilities := []string{"Overgrow"}
	if base.Number == "004" {
		abilities = []string{"Blaze"}
	}

	return &domain.PokemonDetail{
		ID:           base.ID,
		Name:         base.Name,
		Number:       base.Number,
		Types:        mapTypes(base.Types),
		Description:  base.Description,
		ImageURL:     base.ImageURL,
		ElementColor: base.ElementColor,
		Height:       base.Height,
		Weight:       base.Weight,
		Category:     "Seed",
		Abilities:    abilities,
		Weaknesses: []domain.Type{
			{Name: "Fogo", Color: "#EE8130"},
		},
		Evolutions: []domain.Evolution{
			{ID: base.ID, Number: base.Number, Name: base.Name, ImageURL: base.ImageURL, Types: mapTypes(base.Types)},
		},
		Region:     "Kanto",
		Generation: "Geração I",
	}, nil
}

func sanitizePage(page, pageSize int) (int, int) {
	if page < 0 {
		page = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func buildPage(items []domain.Pokemon, total int64, page, pageSize int) *domain.PokemonPage {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return &domain.PokemonPage{
		Content:       items,
		TotalElements: total,
		CurrentPage:   page,
		TotalPages:    totalPages,
		HasNext:       page+1 < totalPages,
	}
}

func pageSlice(items []domain.Pokemon, page, pageSize int) *domain.PokemonPage {
	total := int64(len(items))
	start := page * pageSize
	if start > len(items) {
		start = len(items)
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	slice := make([]domain.Pokemon, 0, end-start)
	if end > start {
		slice = append(slice, items[start:end]...)
	}

	return buildPage(slice, total, page, pageSize)
}

func readRows(rows pgx.Rows) ([]domain.Pokemon, error) {
	items := make([]domain.Pokemon, 0)
	for rows.Next() {
		var p domain.Pokemon
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Number,
			&p.Types,
			&p.Height,
			&p.Weight,
			&p.Description,
			&p.ImageURL,
			&p.ElementColor,
			&p.ElementType,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostgresPokemonRepository) AddFavorite(ctx context.Context, userID, pokemonID string) error {
	query := `
		INSERT INTO user_favorites (user_id, pokemon_id, created_at)
		VALUES ($1::UUID, $2::BIGINT, now())
		ON CONFLICT (user_id, pokemon_id) DO NOTHING
	`
	ct, err := r.db.Exec(ctx, query, userID, pokemonID)
	if err != nil {
		return fmt.Errorf("erro ao adicionar favorito: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrFavoriteAlreadyExists
	}
	return nil
}

func (r *PostgresPokemonRepository) RemoveFavorite(ctx context.Context, userID, pokemonID string) error {
	query := `
		DELETE FROM user_favorites
		WHERE user_id = $1::UUID AND pokemon_id = $2::BIGINT
	`
	ct, err := r.db.Exec(ctx, query, userID, pokemonID)
	if err != nil {
		return fmt.Errorf("erro ao remover favorito: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrFavoriteNotFound
	}
	return nil
}

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, fmt.Errorf("database url vazio")
	}

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 20
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

type evolutionNode struct {
	Pokemon struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"pokemon"`
	Condition struct {
		Description string `json:"description"`
	} `json:"condition"`
	EvolutionsTo []evolutionNode `json:"evolutions_to"`
}

type evolutionStep struct {
	ID      int64
	Name    string
	Trigger string
}

func (r *PostgresPokemonRepository) listAbilityNames(ctx context.Context, id string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT a.name
		FROM pokemons p
		JOIN pokemon_abilities pa ON pa.pokemon_id = p.id
		JOIN abilities a ON a.id = pa.ability_id
		WHERE p.id::text = $1 OR p.number = $1
		ORDER BY a.name ASC
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	abilities := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		abilities = append(abilities, name)
	}

	return abilities, rows.Err()
}

func (r *PostgresPokemonRepository) listWeaknesses(ctx context.Context, id string) ([]domain.Type, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.name, t.color
		FROM pokemons p
		JOIN pokemon_weaknesses pw ON pw.pokemon_id = p.id
		JOIN types t ON t.id = pw.type_id
		WHERE p.id::text = $1 OR p.number = $1
		GROUP BY t.id, t.name, t.color
		ORDER BY t.id ASC
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	weaknesses := make([]domain.Type, 0)
	for rows.Next() {
		var item domain.Type
		if err := rows.Scan(&item.Name, &item.Color); err != nil {
			return nil, err
		}
		weaknesses = append(weaknesses, item)
	}

	return weaknesses, rows.Err()
}

func (r *PostgresPokemonRepository) listEvolutions(ctx context.Context, number string, evolutionChainID int64) ([]domain.Evolution, error) {
	if override, ok := pokemonEvolutionOverrides[strings.TrimLeft(number, "0")]; ok {
		return r.buildEvolutionOverrides(ctx, override), nil
	}

	if evolutionChainID == 0 {
		return []domain.Evolution{}, nil
	}

	var chainRaw []byte
	if err := r.db.QueryRow(ctx, `
		SELECT chain_data
		FROM evolution_chains
		WHERE id = $1
	`, evolutionChainID).Scan(&chainRaw); err != nil {
		return nil, err
	}

	var root evolutionNode
	if err := json.Unmarshal(chainRaw, &root); err != nil {
		return nil, err
	}

	steps := make([]evolutionStep, 0)
	flattenEvolutionChain(root, "", &steps)
	if len(steps) == 0 {
		return []domain.Evolution{}, nil
	}

	evolutions := make([]domain.Evolution, 0, len(steps))
	for _, step := range steps {
		pokemon, err := r.GetByID(ctx, fmt.Sprintf("%d", step.ID))
		if err != nil {
			continue
		}

		evolutions = append(evolutions, domain.Evolution{
			ID:       pokemon.ID,
			Number:   pokemon.Number,
			Name:     pokemon.Name,
			ImageURL: pokemon.ImageURL,
			Types:    mapTypes(pokemon.Types),
			Trigger:  normalizeCondition(step.Trigger),
		})
	}

	return evolutions, nil
}

func (r *PostgresPokemonRepository) buildEvolutionOverrides(ctx context.Context, items []evolutionOverride) []domain.Evolution {
	evolutions := make([]domain.Evolution, 0, len(items))
	for _, item := range items {
		pokemon, err := r.GetByID(ctx, item.Number)
		if err == nil {
			evolutions = append(evolutions, domain.Evolution{
				ID:       pokemon.ID,
				Number:   pokemon.Number,
				Name:     pokemon.Name,
				ImageURL: pokemon.ImageURL,
				Types:    mapTypesWithOverrides(item.Number, pokemon.Types),
				Trigger:  item.Trigger,
			})
			continue
		}

		evolutions = append(evolutions, domain.Evolution{
			ID:       item.Number,
			Number:   item.Number,
			Name:     item.Name,
			ImageURL: fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/%s.png", strings.TrimLeft(item.Number, "0")),
			Types:    mapTypesWithOverrides(item.Number, item.Types),
			Trigger:  item.Trigger,
		})
	}
	return evolutions
}

func flattenEvolutionChain(node evolutionNode, trigger string, out *[]evolutionStep) {
	*out = append(*out, evolutionStep{
		ID:      node.Pokemon.ID,
		Name:    node.Pokemon.Name,
		Trigger: trigger,
	})

	for _, next := range node.EvolutionsTo {
		flattenEvolutionChain(next, next.Condition.Description, out)
	}
}

func mapTypes(typeNames []string) []domain.Type {
	return mapTypesWithOverrides("", typeNames)
}

func mapTypesWithOverrides(number string, typeNames []string) []domain.Type {
	if override, ok := pokemonTypeOverrides[strings.TrimLeft(number, "0")]; ok {
		typeNames = override
	}

	types := make([]domain.Type, len(typeNames))
	for i, name := range typeNames {
		displayName := normalizeTypeName(name)
		types[i] = domain.Type{Name: displayName, Color: typeColor(displayName)}
	}
	return types
}

func normalizeTypeName(value string) string {
	switch strings.TrimSpace(value) {
	case "Aço":
		return "Metal"
	case "Sombrio":
		return "Noturno"
	default:
		return strings.TrimSpace(value)
	}
}

func normalizeCategory(number string, value string) string {
	if override, ok := pokemonCategoryOverrides[strings.TrimLeft(number, "0")]; ok {
		return override
	}

	value = strings.TrimSpace(strings.TrimSuffix(value, " Pokémon"))
	if value == "" {
		return ""
	}
	return value
}

func normalizeCondition(value string) string {
	value = strings.TrimSpace(value)

	switch value {
	case "Subir nivel com felicidade":
		return "Nível de Amizade"
	case "Uso de thunder stone":
		return "Pedra do Trovão"
	case "Uso de moon stone":
		return "Pedra da Lua"
	case "Uso de dusk stone":
		return "Pedra do Anoitecer"
	case "Troca":
		return "Trocas"
	case "Condicao nao mapeada":
		return "Subir de Nível c/ Rollout"
	}

	if strings.HasPrefix(value, "Nivel ") {
		level := strings.TrimPrefix(value, "Nivel ")
		if level == "32" {
			level = "36"
		}
		return "Nível " + level
	}
	return value
}

func normalizeDescription(number string, value string) string {
	if override, ok := pokemonDescriptionOverrides[strings.TrimLeft(number, "0")]; ok {
		return override
	}
	return strings.TrimSpace(value)
}

func generationLabel(id int) string {
	return fmt.Sprintf("%dº Geração", id)
}

func normalizeAbilities(number string, abilities []string) []string {
	if override, ok := pokemonAbilityOverrides[strings.TrimLeft(number, "0")]; ok {
		return override
	}
	return abilities
}

func typeColor(typeName string) string {
	return domain.TypeColor(typeName)
}
