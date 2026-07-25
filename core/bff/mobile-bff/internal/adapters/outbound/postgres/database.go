package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database mantém o pool de conexão com o PostgreSQL.
type Database struct {
	Pool *pgxpool.Pool
}

// NewDatabase cria uma nova conexão com o banco de dados.
func NewDatabase(ctx context.Context, databaseURL string) (*Database, error) {
	if databaseURL == "" {
		return nil, nil
	}

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}
	cfg.MaxConns = 20
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{Pool: pool}, nil
}

// Close fecha o pool de conexões.
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
