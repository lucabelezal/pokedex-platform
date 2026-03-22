package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"pokedex-platform/app/pokedex-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPokemonNotFound = errors.New("pokemon nao encontrado")

type PokemonRepository interface {
	GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error)
	Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error)
	GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error)
	GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
}

type PostgresPokemonRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPokemonRepository(db *pgxpool.Pool) *PostgresPokemonRepository {
	return &PostgresPokemonRepository{db: db}
}

func (r *PostgresPokemonRepository) GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error) {
	page, pageSize = sanitizePage(page, pageSize)
	offset := page * pageSize

	query := `
		SELECT id, name, number, types, height, weight, description, image_url, element_color, element_type, created_at, updated_at
		FROM pokemons
		ORDER BY number ASC
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
		SELECT id, name, number, types, height, weight, description, image_url, element_color, element_type, created_at, updated_at
		FROM pokemons
		WHERE name ILIKE $1 OR number ILIKE $1
		ORDER BY number ASC
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
		SELECT id, name, number, types, height, weight, description, image_url, element_color, element_type, created_at, updated_at
		FROM pokemons
		WHERE $1 = ANY(types)
		ORDER BY number ASC
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
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pokemons WHERE $1 = ANY(types)", typeFilter).Scan(&total); err != nil {
		return nil, err
	}

	return buildPage(items, total, page, pageSize), nil
}

func (r *PostgresPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
	query := `
		SELECT id, name, number, types, height, weight, description, image_url, element_color, element_type, created_at, updated_at
		FROM pokemons
		WHERE id = $1 OR number = $1
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

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, fmt.Errorf("database url vazio")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
