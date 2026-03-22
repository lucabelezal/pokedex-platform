package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"pokedex-platform/app/auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciais invalidas")
	ErrInvalidInput       = errors.New("dados invalidos")
)

type AuthService struct {
	repo      *repository.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

type AuthResult struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int64  `json:"expiresIn"`
	UserID      string `json:"userId"`
	Email       string `json:"email"`
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string, tokenTTLmins int) *AuthService {
	if tokenTTLmins <= 0 {
		tokenTTLmins = 15
	}
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		tokenTTL:  time.Duration(tokenTTLmins) * time.Minute,
	}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)
	if email == "" || password == "" || len(password) < 6 {
		return nil, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, email, string(hash))
	if err != nil {
		return nil, err
	}

	return s.createToken(user.ID, user.Email)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)
	if email == "" || password == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.createToken(user.ID, user.Email)
}

func (s *AuthService) createToken(userID, email string) (*AuthResult, error) {
	now := time.Now()
	expiresAt := now.Add(s.tokenTTL)

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iat":   now.Unix(),
		"exp":   expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.tokenTTL.Seconds()),
		UserID:      userID,
		Email:       email,
	}, nil
}
