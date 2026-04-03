package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresPokemonRepository implementa PokemonRepository usando PostgreSQL.
type PostgresPokemonRepository struct {
	db *pgxpool.Pool
}

// NewPostgresPokemonRepository cria um novo repositório PostgreSQL.
func NewPostgresPokemonRepository(db *pgxpool.Pool) *PostgresPokemonRepository {
	return &PostgresPokemonRepository{db: db}
}

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

func (r *PostgresPokemonRepository) GetDetailByID(ctx context.Context, id string) (*domain.PokemonScreenDetail, error) {
	pokemon, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.PokemonScreenDetail{
		ID:           pokemon.ID,
		Name:         pokemon.Name,
		Number:       pokemon.Number,
		Types:        postgresDomainTypes(pokemon.Types),
		Description:  pokemon.Description,
		ImageURL:     pokemon.ImageURL,
		ElementColor: pokemon.ElementColor,
		Height:       pokemon.Height,
		Weight:       pokemon.Weight,
		Category:     "",
		Abilities:    []string{},
		Weaknesses:   []domain.Type{},
		Evolutions: []domain.Evolution{
			{ID: pokemon.ID, Number: pokemon.Number, Name: pokemon.Name, ImageURL: pokemon.ImageURL, Types: postgresDomainTypes(pokemon.Types)},
		},
	}, nil
}

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

	var totalElements int64
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pokemons").Scan(&totalElements)
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

	var totalElements int64
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pokemons WHERE $1 = ANY(types)", typeFilter).Scan(&totalElements)
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

func (r *PostgresPokemonRepository) ListTypes(ctx context.Context) ([]domain.Type, error) {
	query := `
		SELECT name, color
		FROM types
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	types := make([]domain.Type, 0)
	for rows.Next() {
		var pokemonType domain.Type
		if err := rows.Scan(&pokemonType.Name, &pokemonType.Color); err != nil {
			return nil, err
		}
		types = append(types, pokemonType)
	}

	return types, rows.Err()
}

func (r *PostgresPokemonRepository) ListRegions(ctx context.Context) ([]domain.Region, error) {
	rows, err := r.db.Query(ctx, `
		SELECT name
		FROM regions
		WHERE id <= 8
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	regions := make([]domain.Region, 0)
	index := 1
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		regions = append(regions, domain.Region{
			ID:         strings.ToLower(name),
			Name:       name,
			Generation: fmt.Sprintf("%dº Geração", index),
		})
		index++
	}

	return regions, rows.Err()
}

func postgresDomainTypes(names []string) []domain.Type {
	result := make([]domain.Type, len(names))
	for i, name := range names {
		result[i] = domain.Type{Name: name, Color: postgresTypeColor(name)}
	}
	return result
}

func postgresTypeColor(name string) string {
	switch name {
	case "Normal":
		return "#A8A77A"
	case "Fogo":
		return "#EE8130"
	case "Água":
		return "#6390F0"
	case "Elétrico":
		return "#F7D02C"
	case "Grama":
		return "#7AC74C"
	case "Gelo":
		return "#96D9D6"
	case "Lutador":
		return "#C22E28"
	case "Venenoso":
		return "#A33EA1"
	case "Terrestre":
		return "#E2BF65"
	case "Voador":
		return "#A98FF3"
	case "Psíquico":
		return "#F95587"
	case "Inseto":
		return "#A6B91A"
	case "Pedra":
		return "#B6A136"
	case "Fantasma":
		return "#735797"
	case "Dragão":
		return "#6F35FC"
	case "Sombrio", "Noturno":
		return "#705746"
	case "Aço", "Metal":
		return "#B7B7CE"
	case "Fada":
		return "#D685AD"
	default:
		return "#A9AC86"
	}
}
