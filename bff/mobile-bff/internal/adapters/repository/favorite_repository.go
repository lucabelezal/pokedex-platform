package repository

import (
	"context"
	"errors"
	"time"

	"pokedex-platform/bff/mobile-bff/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresFavoriteRepository implements the FavoriteRepository port using PostgreSQL
type PostgresFavoriteRepository struct {
	db *pgxpool.Pool
}

// NewPostgresFavoriteRepository creates a new PostgreSQL favorites repository
func NewPostgresFavoriteRepository(db *pgxpool.Pool) *PostgresFavoriteRepository {
	return &PostgresFavoriteRepository{db: db}
}

// AddFavorite adds a Pokemon to user's favorites
func (r *PostgresFavoriteRepository) AddFavorite(ctx context.Context, userID, pokemonID string) error {
	// Check if already favorited
	isFav, err := r.IsFavorite(ctx, userID, pokemonID)
	if err != nil {
		return err
	}
	if isFav {
		return domain.ErrFavoriteAlreadyExists
	}

	query := `
		INSERT INTO favorites (id, user_id, pokemon_id, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3)
	`

	_, err = r.db.Exec(ctx, query, userID, pokemonID, time.Now())
	return err
}

// RemoveFavorite removes a Pokemon from user's favorites
func (r *PostgresFavoriteRepository) RemoveFavorite(ctx context.Context, userID, pokemonID string) error {
	query := `
		DELETE FROM favorites
		WHERE user_id = $1 AND pokemon_id = $2
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

// IsFavorite checks if a Pokemon is in user's favorites
func (r *PostgresFavoriteRepository) IsFavorite(ctx context.Context, userID, pokemonID string) (bool, error) {
	query := `
		SELECT 1 FROM favorites
		WHERE user_id = $1 AND pokemon_id = $2
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

// GetUserFavorites retrieves all favorite Pokemon IDs for a user
func (r *PostgresFavoriteRepository) GetUserFavorites(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return []string{}, nil
	}

	query := `
		SELECT pokemon_id
		FROM favorites
		WHERE user_id = $1
		ORDER BY created_at DESC
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
