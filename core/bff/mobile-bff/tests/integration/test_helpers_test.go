package integration

import (
	"context"
	"testing"
	"time"

	repository "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/postgres"
)

const testUserID = "11111111-1111-1111-1111-111111111111"

// seedTestUser garante a existência de um usuário com UUID válido (o repo usa $1::UUID).
func seedTestUser(t *testing.T, db *repository.Database) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Pool.Exec(ctx, `
		INSERT INTO users (id, email) VALUES ($1, 'escrita@local')
		ON CONFLICT (id) DO NOTHING
	`, testUserID)
	if err != nil {
		t.Fatalf("falha ao inserir usuario de escrita: %v", err)
	}
}
