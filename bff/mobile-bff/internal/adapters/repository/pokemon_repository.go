package repository

import (
	"context"
	"errors"

	"pokedex-platform/bff/mobile-bff/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresPokemonRepository implements the PokemonRepository port using PostgreSQL
type PostgresPokemonRepository struct {
	db *pgxpool.Pool
}

// NewPostgresPokemonRepository creates a new PostgreSQL repository
func NewPostgresPokemonRepository(db *pgxpool.Pool) *PostgresPokemonRepository {
	return &PostgresPokemonRepository{db: db}
}

// GetByID retrieves a single Pokemon by ID
func (r *PostgresPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
	query := `
		SELECT id, name, number, types, height, weight, description, image_url, element_color, element_type, created_at, updated_at
		FROM pokemons
		WHERE id = $1
	`

	var pokemon domain.Pokemon
	err := r.db.QueryRow(ctx, query, id).Scan(
		&pokemon.ID,
		&pokemon.Name,
		&pokemon.Number,
		&pokemon.Types,
		&pokemon.Height,
		&pokemon.Weight,
		&pokemon.Description,
		&pokemon.ImageURL,
		&pokemon.ElementColor,
		&pokemon.ElementType,
		&pokemon.CreatedAt,
		&pokemon.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		return nil, domain.ErrPokemonNotFound
	}

	return &pokemon, nil
}

// GetAll retrieves all Pokemons with pagination
func (r *PostgresPokemonRepository) GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page < 0 {
		page = 0
	}

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

	var pokemons []domain.Pokemon
	for rows.Next() {
		var pokemon domain.Pokemon
		err := rows.Scan(
			&pokemon.ID,
			&pokemon.Name,
			&pokemon.Number,
			&pokemon.Types,
			&pokemon.Height,
			&pokemon.Weight,
			&pokemon.Description,
			&pokemon.ImageURL,
			&pokemon.ElementColor,
			&pokemon.ElementType,
			&pokemon.CreatedAt,
			&pokemon.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pokemons = append(pokemons, pokemon)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Get total count
	var totalElements int64
	countQuery := "SELECT COUNT(*) FROM pokemons"
	err = r.db.QueryRow(ctx, countQuery).Scan(&totalElements)
	if err != nil {
		return nil, err
	}

	totalPages := int((totalElements + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page+1 < totalPages

	return &domain.PokemonPage{
		Content:       pokemons,
		TotalElements: totalElements,
		CurrentPage:   page,
		TotalPages:    totalPages,
		HasNext:       hasNext,
	}, nil
}

// Search retrieves Pokemons matching a query string
func (r *PostgresPokemonRepository) Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page < 0 {
		page = 0
	}

	offset := page * pageSize
	searchQuery := "%" + query + "%"

	sql := `
		SELECT id, name, number, types, height, weight, description, image_url, element_color, element_type, created_at, updated_at
		FROM pokemons
		WHERE name ILIKE $1 OR number ILIKE $1
		ORDER BY number ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, sql, searchQuery, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokemons []domain.Pokemon
	for rows.Next() {
		var pokemon domain.Pokemon
		err := rows.Scan(
			&pokemon.ID,
			&pokemon.Name,
			&pokemon.Number,
			&pokemon.Types,
			&pokemon.Height,
			&pokemon.Weight,
			&pokemon.Description,
			&pokemon.ImageURL,
			&pokemon.ElementColor,
			&pokemon.ElementType,
			&pokemon.CreatedAt,
			&pokemon.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pokemons = append(pokemons, pokemon)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Get total count
	var totalElements int64
	countSQL := "SELECT COUNT(*) FROM pokemons WHERE name ILIKE $1 OR number ILIKE $1"
	err = r.db.QueryRow(ctx, countSQL, searchQuery).Scan(&totalElements)
	if err != nil {
		return nil, err
	}

	totalPages := int((totalElements + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page+1 < totalPages

	return &domain.PokemonPage{
		Content:       pokemons,
		TotalElements: totalElements,
		CurrentPage:   page,
		TotalPages:    totalPages,
		HasNext:       hasNext,
	}, nil
}

// GetByType retrieves Pokemons filtered by type
func (r *PostgresPokemonRepository) GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page < 0 {
		page = 0
	}

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

	var pokemons []domain.Pokemon
	for rows.Next() {
		var pokemon domain.Pokemon
		err := rows.Scan(
			&pokemon.ID,
			&pokemon.Name,
			&pokemon.Number,
			&pokemon.Types,
			&pokemon.Height,
			&pokemon.Weight,
			&pokemon.Description,
			&pokemon.ImageURL,
			&pokemon.ElementColor,
			&pokemon.ElementType,
			&pokemon.CreatedAt,
			&pokemon.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pokemons = append(pokemons, pokemon)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Get total count
	var totalElements int64
	countQuery := "SELECT COUNT(*) FROM pokemons WHERE $1 = ANY(types)"
	err = r.db.QueryRow(ctx, countQuery, typeFilter).Scan(&totalElements)
	if err != nil {
		return nil, err
	}

	totalPages := int((totalElements + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page+1 < totalPages

	return &domain.PokemonPage{
		Content:       pokemons,
		TotalElements: totalElements,
		CurrentPage:   page,
		TotalPages:    totalPages,
		HasNext:       hasNext,
	}, nil
}

// GetFavorites retrieves favorite Pokemon IDs for a user
func (r *PostgresPokemonRepository) GetFavorites(ctx context.Context, userID string, page, pageSize int) ([]string, error) {
	if userID == "" {
		return []string{}, nil
	}

	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page < 0 {
		page = 0
	}

	offset := page * pageSize

	query := `
		SELECT pokemon_id
		FROM favorites
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []string
	for rows.Next() {
		var pokemonID string
		if err := rows.Scan(&pokemonID); err != nil {
			return nil, err
		}
		favorites = append(favorites, pokemonID)
	}

	return favorites, rows.Err()
}
