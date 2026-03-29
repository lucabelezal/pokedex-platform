package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound         = errors.New("usuario nao encontrado")
	ErrUserAlreadyExist     = errors.New("usuario ja existe")
	ErrRefreshTokenNotFound = errors.New("refresh token nao encontrado")
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type RefreshSession struct {
	UserID    string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
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

func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash string) (*User, error) {
	query := `
		INSERT INTO auth_users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at
	`

	user := &User{}
	err := r.db.QueryRow(ctx, query, strings.ToLower(strings.TrimSpace(email)), passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, ErrUserAlreadyExist
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM auth_users
		WHERE email = $1
	`

	user := &User{}
	err := r.db.QueryRow(ctx, query, strings.ToLower(strings.TrimSpace(email))).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM auth_users
		WHERE id = $1::UUID
	`

	user := &User{}
	err := r.db.QueryRow(ctx, query, strings.TrimSpace(userID)).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) StoreRefreshToken(ctx context.Context, userID, refreshToken string, expiresAt time.Time) error {
	query := `
		INSERT INTO auth_refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1::UUID, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, strings.TrimSpace(userID), hashRefreshToken(refreshToken), expiresAt.UTC())
	return err
}

func (r *UserRepository) GetActiveRefreshSession(ctx context.Context, refreshToken string) (*RefreshSession, error) {
	query := `
		SELECT user_id, expires_at, revoked_at
		FROM auth_refresh_tokens
		WHERE token_hash = $1
	`

	session := &RefreshSession{}
	err := r.db.QueryRow(ctx, query, hashRefreshToken(refreshToken)).Scan(&session.UserID, &session.ExpiresAt, &session.RevokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	if session.RevokedAt != nil || time.Now().UTC().After(session.ExpiresAt.UTC()) {
		return nil, ErrRefreshTokenNotFound
	}

	return session, nil
}

func (r *UserRepository) RotateRefreshToken(ctx context.Context, currentToken, newToken, userID string, expiresAt time.Time) error {
	currentHash := hashRefreshToken(currentToken)
	newHash := hashRefreshToken(newToken)
	now := time.Now().UTC()

	return pgx.BeginFunc(ctx, r.db, func(tx pgx.Tx) error {
		var currentUserID string
		var currentExpiresAt time.Time
		var revokedAt *time.Time

		query := `
			SELECT user_id, expires_at, revoked_at
			FROM auth_refresh_tokens
			WHERE token_hash = $1
			FOR UPDATE
		`
		if err := tx.QueryRow(ctx, query, currentHash).Scan(&currentUserID, &currentExpiresAt, &revokedAt); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrRefreshTokenNotFound
			}
			return err
		}

		if currentUserID != strings.TrimSpace(userID) || revokedAt != nil || now.After(currentExpiresAt.UTC()) {
			return ErrRefreshTokenNotFound
		}

		if _, err := tx.Exec(ctx, `
			UPDATE auth_refresh_tokens
			SET revoked_at = $2, replaced_by_token_hash = $3
			WHERE token_hash = $1
		`, currentHash, now, newHash); err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO auth_refresh_tokens (user_id, token_hash, expires_at)
			VALUES ($1::UUID, $2, $3)
		`, currentUserID, newHash, expiresAt.UTC()); err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	query := `
		UPDATE auth_refresh_tokens
		SET revoked_at = COALESCE(revoked_at, $2)
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > $2
	`

	result, err := r.db.Exec(ctx, query, hashRefreshToken(refreshToken), time.Now().UTC())
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

func (r *UserRepository) RevokeAccessToken(ctx context.Context, jti string, expiresAt time.Time) error {
	jti = strings.TrimSpace(jti)
	if jti == "" {
		return nil
	}

	query := `
		INSERT INTO auth_revoked_access_tokens (jti, expires_at, revoked_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (jti)
		DO UPDATE SET
			expires_at = GREATEST(auth_revoked_access_tokens.expires_at, EXCLUDED.expires_at),
			revoked_at = EXCLUDED.revoked_at
	`

	_, err := r.db.Exec(ctx, query, jti, expiresAt.UTC(), time.Now().UTC())
	return err
}

func (r *UserRepository) IsAccessTokenRevoked(ctx context.Context, jti string) (bool, error) {
	jti = strings.TrimSpace(jti)
	if jti == "" {
		return true, nil
	}

	query := `
		SELECT 1
		FROM auth_revoked_access_tokens
		WHERE jti = $1 AND expires_at > $2
		LIMIT 1
	`

	var exists int
	err := r.db.QueryRow(ctx, query, jti, time.Now().UTC()).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *UserRepository) CleanupExpiredAuthArtifacts(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM auth_revoked_access_tokens
		WHERE expires_at <= $1
	`, time.Now().UTC())
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		DELETE FROM auth_refresh_tokens
		WHERE expires_at <= $1
	`, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}

func hashRefreshToken(refreshToken string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(refreshToken)))
	return hex.EncodeToString(sum[:])
}
