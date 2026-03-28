package repository

import (
	"context"
	"errors"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresFavoriteRepository implementa a interface FavoriteRepository usando PostgreSQL
type PostgresFavoriteRepository struct {
	db *pgxpool.Pool
}

// NewPostgresFavoriteRepository cria um novo repositório PostgreSQL de favoritos
func NewPostgresFavoriteRepository(db *pgxpool.Pool) *PostgresFavoriteRepository {
	return &PostgresFavoriteRepository{db: db}
}

// AddFavorite adiciona um Pokémon aos favoritos do usuário
func (r *PostgresFavoriteRepository) AddFavorite(ctx context.Context, userID, pokemonID string) error {
	isFav, err := r.IsFavorite(ctx, userID, pokemonID)
	if err != nil {
		return err
	}
	if isFav {
		return domain.ErrFavoriteAlreadyExists
	}

	query := `
		INSERT INTO user_favorites (user_id, pokemon_id, created_at)
		VALUES ($1::UUID, $2, $3)
		ON CONFLICT (user_id, pokemon_id) DO NOTHING
	`
	_, err = r.db.Exec(ctx, query, userID, pokemonID, time.Now())
	return err
}

// RemoveFavorite remove um Pokémon dos favoritos do usuário
func (r *PostgresFavoriteRepository) RemoveFavorite(ctx context.Context, userID, pokemonID string) error {
	query := `
		DELETE FROM user_favorites
		WHERE user_id = $1::UUID AND pokemon_id = $2
	`

	result, err := r.db.Exec(ctx, query, userID, pokemonID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrFavoriteNotFound
	}

	return nil
}

// IsFavorite verifica se um Pokémon está nos favoritos do usuário
func (r *PostgresFavoriteRepository) IsFavorite(ctx context.Context, userID, pokemonID string) (bool, error) {
	query := `
		SELECT 1 FROM user_favorites
		WHERE user_id = $1::UUID AND pokemon_id = $2
		LIMIT 1
	`

	var exists int
	err := r.db.QueryRow(ctx, query, userID, pokemonID).Scan(&exists)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetUserFavorites recupera todos os IDs de Pokémons favoritos do usuário
func (r *PostgresFavoriteRepository) GetUserFavorites(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return []string{}, nil
	}

	if !isValidUUID(userID) {
		return []string{}, nil
	}

	query := `
		SELECT DISTINCT pokemon_id::TEXT
		FROM user_favorites
		WHERE user_id = $1::UUID
	`

	rows, err := r.db.Query(ctx, query, userID)
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

// isValidUUID verifica se uma string é um UUID válido no formato xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (36 caracteres com hífens)
func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	// Verificar posições de hífens (formato 8-4-4-4-12)
	return s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-'
}
