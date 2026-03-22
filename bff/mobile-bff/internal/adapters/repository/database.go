package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database mantém o pool de conexão
type Database struct {
	Pool *pgxpool.Pool
}

// NewDatabase cria uma nova conexão com o banco de dados
func NewDatabase(ctx context.Context, databaseURL string) (*Database, error) {
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5432/pokedex"
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Testar a conexão
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{Pool: pool}, nil
}

// Close fecha a conexão com o banco de dados
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
