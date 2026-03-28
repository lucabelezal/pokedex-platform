package ports

import "context"

// AuthSession represents the authentication payload returned by an auth provider.
type AuthSession struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
	UserID      string
	Email       string
}

// AuthProvider defines the outbound authentication operations needed by the BFF.
type AuthProvider interface {
	Signup(ctx context.Context, email, password string) (*AuthSession, error)
	Login(ctx context.Context, email, password string) (*AuthSession, error)
	Refresh(ctx context.Context, token string) (*AuthSession, error)
	Logout(ctx context.Context, token string) error
}
