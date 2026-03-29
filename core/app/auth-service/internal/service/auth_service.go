package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"pokedex-platform/core/app/auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciais invalidas")
	ErrInvalidInput       = errors.New("dados invalidos")
	ErrInvalidToken       = errors.New("token invalido")
)

const (
	minPasswordLength = 8
	bcryptCost        = 12
)

type AuthService struct {
	repo            AuthRepository
	jwtSecret       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type AuthRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*repository.User, error)
	GetByEmail(ctx context.Context, email string) (*repository.User, error)
	GetByID(ctx context.Context, userID string) (*repository.User, error)
	StoreRefreshToken(ctx context.Context, userID, refreshToken string, expiresAt time.Time) error
	GetActiveRefreshSession(ctx context.Context, refreshToken string) (*repository.RefreshSession, error)
	RotateRefreshToken(ctx context.Context, currentToken, newToken, userID string, expiresAt time.Time) error
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	RevokeAccessToken(ctx context.Context, jti string, expiresAt time.Time) error
	IsAccessTokenRevoked(ctx context.Context, jti string) (bool, error)
}

type AuthResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
	UserID       string `json:"userId"`
	Email        string `json:"email"`
}

func NewAuthService(repo AuthRepository, jwtSecret string, tokenTTLmins int, refreshTokenTTLHours int) *AuthService {
	if tokenTTLmins <= 0 {
		tokenTTLmins = 15
	}
	if refreshTokenTTLHours <= 0 {
		refreshTokenTTLHours = 168
	}
	return &AuthService{
		repo:            repo,
		jwtSecret:       jwtSecret,
		accessTokenTTL:  time.Duration(tokenTTLmins) * time.Minute,
		refreshTokenTTL: time.Duration(refreshTokenTTLHours) * time.Hour,
	}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)
	if !isValidEmail(email) || len(password) < minPasswordLength {
		return nil, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, email, string(hash))
	if err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user.ID, user.Email)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)
	if !isValidEmail(email) || password == "" {
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

	return s.issueSession(ctx, user.ID, user.Email)
}

func (s *AuthService) Refresh(ctx context.Context, tokenString string) (*AuthResult, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	session, err := s.repo.GetActiveRefreshSession(ctx, tokenString)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	user, userErr := s.repo.GetByID(ctx, session.UserID)
	if userErr != nil {
		return nil, ErrInvalidToken
	}

	return s.rotateRefreshSession(ctx, tokenString, user.ID, user.Email)
}

func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return ErrInvalidToken
	}

	if err := s.repo.RevokeRefreshToken(ctx, tokenString); err == nil {
		return nil
	} else if !errors.Is(err, repository.ErrRefreshTokenNotFound) {
		return err
	}

	claims, err := s.parseAndValidateToken(tokenString)
	if err != nil {
		return ErrInvalidToken
	}

	jti, expiresAt, err := extractJTIAndExpiresAt(claims)
	if err != nil {
		return ErrInvalidToken
	}

	if err := s.repo.RevokeAccessToken(ctx, jti, expiresAt); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) IsAccessTokenActive(ctx context.Context, tokenString string) (bool, error) {
	claims, err := s.parseAndValidateToken(tokenString)
	if err != nil {
		return false, nil
	}

	jti, _, err := extractJTIAndExpiresAt(claims)
	if err != nil {
		return false, nil
	}

	revoked, err := s.repo.IsAccessTokenRevoked(ctx, jti)
	if err != nil {
		return false, err
	}

	return !revoked, nil
}

func (s *AuthService) issueSession(ctx context.Context, userID, email string) (*AuthResult, error) {
	refreshToken, refreshExpiresAt, err := generateRefreshToken(s.refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	if err := s.repo.StoreRefreshToken(ctx, userID, refreshToken, refreshExpiresAt); err != nil {
		return nil, err
	}

	return s.createAuthResult(userID, email, refreshToken)
}

func (s *AuthService) rotateRefreshSession(ctx context.Context, currentRefreshToken, userID, email string) (*AuthResult, error) {
	newRefreshToken, refreshExpiresAt, err := generateRefreshToken(s.refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	if err := s.repo.RotateRefreshToken(ctx, currentRefreshToken, newRefreshToken, userID, refreshExpiresAt); err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	return s.createAuthResult(userID, email, newRefreshToken)
}

func (s *AuthService) createAuthResult(userID, email, refreshToken string) (*AuthResult, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessTokenTTL)
	tokenID, err := generateTokenID()
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iat":   now.Unix(),
		"exp":   expiresAt.Unix(),
		"jti":   tokenID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		AccessToken:  signed,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
		UserID:       userID,
		Email:        email,
	}, nil
}

func (s *AuthService) parseAndValidateToken(tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de assinatura invalido: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if expRaw, ok := claims["exp"]; ok {
		exp, ok := expRaw.(float64)
		if !ok {
			return nil, ErrInvalidToken
		}
		if time.Now().Unix() >= int64(exp) {
			return nil, ErrInvalidToken
		}
	}

	return claims, nil
}

func isValidEmail(email string) bool {
	parsed, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	return strings.EqualFold(strings.TrimSpace(parsed.Address), email)
}

func generateRefreshToken(ttl time.Duration) (string, time.Time, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", time.Time{}, err
	}

	return base64.RawURLEncoding.EncodeToString(buf), time.Now().UTC().Add(ttl), nil
}

func generateTokenID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func extractJTIAndExpiresAt(claims jwt.MapClaims) (string, time.Time, error) {
	jti, _ := claims["jti"].(string)
	jti = strings.TrimSpace(jti)
	if jti == "" {
		return "", time.Time{}, ErrInvalidToken
	}

	expRaw, ok := claims["exp"]
	if !ok {
		return "", time.Time{}, ErrInvalidToken
	}

	expFloat, ok := expRaw.(float64)
	if !ok {
		return "", time.Time{}, ErrInvalidToken
	}

	return jti, time.Unix(int64(expFloat), 0).UTC(), nil
}
