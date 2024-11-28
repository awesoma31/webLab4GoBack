package service

import (
	"awesoma31/common/api"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtUser struct {
	ID           int64
	Username     string
	AccessToken  string
	RefreshToken string
}

type TokenService interface {
	GenerateToken(id int64, username string) (string, error)
	GenerateTokens(id int64, username string) (api.Tokens, error)
	GenerateRefreshToken(id int64, username string) (string, error)
	ParseToken(ctx context.Context, tokenString string) (*JwtUser, error)
	ValidateToken(tokenString string) bool
	GetUsernameFromToken(tokenString string) (string, error)
	GetUserIDFromToken(tokenString string) (int64, error)
}

type tokenServiceImpl struct {
	secretKey           []byte
	jwtExpirationMs     time.Duration
	refreshExpirationMs time.Duration
}

func NewTokenService(secretKey string, jwtExpirationMs, refreshExpirationMs time.Duration) TokenService {
	return &tokenServiceImpl{
		secretKey:           []byte(secretKey),
		jwtExpirationMs:     jwtExpirationMs,
		refreshExpirationMs: refreshExpirationMs,
	}
}

func (ts *tokenServiceImpl) GenerateTokens(id int64, username string) (api.Tokens, error) {
	at, err := ts.GenerateToken(id, username)
	if err != nil {
		return api.Tokens{}, err
	}
	rt, err := ts.GenerateRefreshToken(id, username)
	if err != nil {
		return api.Tokens{}, err
	}

	return api.Tokens{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (ts *tokenServiceImpl) GenerateToken(id int64, username string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      now.Add(ts.jwtExpirationMs).Unix(),
		"iat":      now.Unix(),
		"iss":      "web-lab4",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secretKey)
}

func (ts *tokenServiceImpl) GenerateRefreshToken(id int64, username string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      now.Add(ts.refreshExpirationMs).Unix(),
		"iat":      now.Unix(),
		"authorities": []string{
			"ROLE_USER",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secretKey)
}

func (ts *tokenServiceImpl) ParseToken(ctx context.Context, tokenString string) (*JwtUser, error) {
	// Add logic to handle context cancellation if necessary
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ts.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &JwtUser{
			ID:          int64(claims["id"].(float64)),
			Username:    claims["username"].(string),
			AccessToken: tokenString,
		}
		return user, nil
	}

	return nil, errors.New("invalid token")
}

func (ts *tokenServiceImpl) ValidateToken(tokenString string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := ts.ParseToken(ctx, tokenString)
	return err == nil
}

func (ts *tokenServiceImpl) GetUsernameFromToken(tokenString string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	user, err := ts.ParseToken(ctx, tokenString)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func (ts *tokenServiceImpl) GetUserIDFromToken(tokenString string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	user, err := ts.ParseToken(ctx, tokenString)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
