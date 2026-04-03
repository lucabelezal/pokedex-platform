package domain

// AuthSession representa o payload de autenticação retornado pelo provedor de auth.
type AuthSession struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
	UserID       string
	Email        string
}
