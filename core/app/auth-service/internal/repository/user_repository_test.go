package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustRepoMock(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	return mock
}

func TestUserRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(pgxmock.PgxPoolIface)
		wantErr error
	}{
		{
			name: "sucesso",
			setup: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta("INSERT INTO auth_users")).
					WithArgs("ash@kanto.dev", "hash").
					WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "created_at"}).
						AddRow("user-1", "ash@kanto.dev", "hash", time.Now()))
			},
		},
		{
			name: "email duplicado",
			setup: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta("INSERT INTO auth_users")).
					WithArgs("ash@kanto.dev", "hash").
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			wantErr: ErrUserAlreadyExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mustRepoMock(t)
			defer func() { mock.Close() }()
			tt.setup(mock)

			r := NewUserRepositoryWithPool(mock)
			_, err := r.CreateUser(context.Background(), "Ash@Kanto.Dev ", "hash")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	t.Run("encontrado", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, created_at")).
			WithArgs("ash@kanto.dev").
			WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "created_at"}).
				AddRow("user-1", "ash@kanto.dev", "hash", time.Now()))

		r := NewUserRepositoryWithPool(mock)
		u, err := r.GetByEmail(context.Background(), "Ash@Kanto.Dev ")
		require.NoError(t, err)
		assert.Equal(t, "user-1", u.ID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("nao encontrado", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, created_at")).
			WithArgs("nobody@x.dev").
			WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "created_at"}))

		r := NewUserRepositoryWithPool(mock)
		_, err := r.GetByEmail(context.Background(), "nobody@x.dev")
		require.ErrorIs(t, err, ErrUserNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	mock := mustRepoMock(t)
	defer func() { mock.Close() }()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, created_at")).
		WithArgs("user-1").
		WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "created_at"}).
			AddRow("user-1", "ash@kanto.dev", "hash", time.Now()))

	r := NewUserRepositoryWithPool(mock)
	u, err := r.GetByID(context.Background(), "user-1")
	require.NoError(t, err)
	assert.Equal(t, "ash@kanto.dev", u.Email)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_StoreRefreshToken(t *testing.T) {
	mock := mustRepoMock(t)
	defer func() { mock.Close() }()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO auth_refresh_tokens")).
		WithArgs("user-1", pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	r := NewUserRepositoryWithPool(mock)
	err := r.StoreRefreshToken(context.Background(), "user-1", "refresh-token", time.Now().Add(time.Hour))
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetActiveRefreshSession(t *testing.T) {
	t.Run("sessao ativa", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, expires_at, revoked_at")).
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "expires_at", "revoked_at"}).
				AddRow("user-1", time.Now().Add(time.Hour), nil))

		r := NewUserRepositoryWithPool(mock)
		s, err := r.GetActiveRefreshSession(context.Background(), "refresh-token")
		require.NoError(t, err)
		assert.Equal(t, "user-1", s.UserID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("nao encontrado", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, expires_at, revoked_at")).
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "expires_at", "revoked_at"}))

		r := NewUserRepositoryWithPool(mock)
		_, err := r.GetActiveRefreshSession(context.Background(), "refresh-token")
		require.ErrorIs(t, err, ErrRefreshTokenNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("sessao revogada", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		revoked := time.Now()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, expires_at, revoked_at")).
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "expires_at", "revoked_at"}).
				AddRow("user-1", time.Now().Add(time.Hour), &revoked))

		r := NewUserRepositoryWithPool(mock)
		_, err := r.GetActiveRefreshSession(context.Background(), "refresh-token")
		require.ErrorIs(t, err, ErrRefreshTokenNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_RevokeRefreshToken(t *testing.T) {
	t.Run("sucesso", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE auth_refresh_tokens")).
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		r := NewUserRepositoryWithPool(mock)
		err := r.RevokeRefreshToken(context.Background(), "refresh-token")
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("nao encontrado", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE auth_refresh_tokens")).
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		r := NewUserRepositoryWithPool(mock)
		err := r.RevokeRefreshToken(context.Background(), "refresh-token")
		require.ErrorIs(t, err, ErrRefreshTokenNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_RevokeAccessToken(t *testing.T) {
	t.Run("jti vazio no-op", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		r := NewUserRepositoryWithPool(mock)
		err := r.RevokeAccessToken(context.Background(), "", time.Now())
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("sucesso", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO auth_revoked_access_tokens")).
			WithArgs("jti-1", pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		r := NewUserRepositoryWithPool(mock)
		err := r.RevokeAccessToken(context.Background(), "jti-1", time.Now().Add(time.Hour))
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_IsAccessTokenRevoked(t *testing.T) {
	t.Run("jti vazio retorna true", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		r := NewUserRepositoryWithPool(mock)
		revoked, err := r.IsAccessTokenRevoked(context.Background(), "")
		require.NoError(t, err)
		assert.True(t, revoked)
	})

	t.Run("revogado", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1")).
			WithArgs("jti-1", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(1))

		r := NewUserRepositoryWithPool(mock)
		revoked, err := r.IsAccessTokenRevoked(context.Background(), "jti-1")
		require.NoError(t, err)
		assert.True(t, revoked)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("nao revogado", func(t *testing.T) {
		mock := mustRepoMock(t)
		defer func() { mock.Close() }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1")).
			WithArgs("jti-1", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}))

		r := NewUserRepositoryWithPool(mock)
		revoked, err := r.IsAccessTokenRevoked(context.Background(), "jti-1")
		require.NoError(t, err)
		assert.False(t, revoked)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_CleanupExpiredAuthArtifacts(t *testing.T) {
	mock := mustRepoMock(t)
	defer func() { mock.Close() }()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM auth_revoked_access_tokens")).
		WithArgs(pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("DELETE", 2))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM auth_refresh_tokens")).
		WithArgs(pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	r := NewUserRepositoryWithPool(mock)
	err := r.CleanupExpiredAuthArtifacts(context.Background())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHashRefreshToken(t *testing.T) {
	got := hashRefreshToken("  refresh-token  ")
	assert.NotEmpty(t, got)
	assert.Len(t, got, 64)
	assert.Equal(t, hashRefreshToken("refresh-token"), got)
}

func TestNewPool_ComURIVazia(t *testing.T) {
	_, err := NewPool(context.Background(), "")
	require.Error(t, err)
}

func TestUserRepository_RotateRefreshToken_Sucesso(t *testing.T) {
	mock := mustRepoMock(t)
	defer func() { mock.Close() }()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, expires_at, revoked_at")).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"user_id", "expires_at", "revoked_at"}).
			AddRow("user-1", time.Now().Add(time.Hour), nil))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE auth_refresh_tokens")).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO auth_refresh_tokens")).
		WithArgs("user-1", pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()
	mock.ExpectRollback()

	r := NewUserRepositoryWithPool(mock)
	err := r.RotateRefreshToken(context.Background(), "old-token", "new-token", "user-1", time.Now().Add(time.Hour))
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_RotateRefreshToken_UsuarioDiferente(t *testing.T) {
	mock := mustRepoMock(t)
	defer func() { mock.Close() }()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, expires_at, revoked_at")).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"user_id", "expires_at", "revoked_at"}).
			AddRow("outro-user", time.Now().Add(time.Hour), nil))
	mock.ExpectRollback().Maybe()
	mock.ExpectRollback().Maybe()

	r := NewUserRepositoryWithPool(mock)
	err := r.RotateRefreshToken(context.Background(), "old-token", "new-token", "user-1", time.Now().Add(time.Hour))
	require.ErrorIs(t, err, ErrRefreshTokenNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_RotateRefreshToken_NaoEncontrado(t *testing.T) {
	mock := mustRepoMock(t)
	defer func() { mock.Close() }()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, expires_at, revoked_at")).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"user_id", "expires_at", "revoked_at"}))
	mock.ExpectRollback().Maybe()
	mock.ExpectRollback().Maybe()

	r := NewUserRepositoryWithPool(mock)
	err := r.RotateRefreshToken(context.Background(), "old-token", "new-token", "user-1", time.Now().Add(time.Hour))
	require.ErrorIs(t, err, ErrRefreshTokenNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}
