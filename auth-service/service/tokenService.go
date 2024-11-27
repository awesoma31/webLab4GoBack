package service

import (
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

type TokenService struct {
	secretKey           []byte
	jwtExpirationMs     time.Duration
	refreshExpirationMs time.Duration
}

func NewTokenService(secretKey string, jwtExpirationMs, refreshExpirationMs time.Duration) *TokenService {
	return &TokenService{
		secretKey:           []byte(secretKey),
		jwtExpirationMs:     jwtExpirationMs,
		refreshExpirationMs: refreshExpirationMs,
	}
}

func (ts *TokenService) GenerateToken(id int64, username string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      now.Add(ts.jwtExpirationMs).Unix(),
		"iat":      now.Unix(),
		"iss":      "your_app_name",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secretKey)
}

func (ts *TokenService) GenerateRefreshToken(id int64, username string) (string, error) {
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

func (ts *TokenService) ParseToken(tokenString string) (*JwtUser, error) {
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

func (ts *TokenService) ValidateToken(tokenString string) bool {
	_, err := ts.ParseToken(tokenString)
	return err == nil
}

func (ts *TokenService) GetUsernameFromToken(tokenString string) (string, error) {
	user, err := ts.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func (ts *TokenService) GetUserIDFromToken(tokenString string) (int64, error) {
	user, err := ts.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
