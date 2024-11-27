package service

import (
	"awesoma31/common/api"
	"context"
	"github.com/awesoma31/auth-service/storage"
	"time"
)

type AuthService interface {
	Authorize(ctx context.Context, r *api.AuthorizeRequest) (*api.Authorization, error)
}

type Service struct {
	store        storage.UsersStore
	tokenService *TokenService
}

func NewAuthService(store storage.UsersStore) *Service {
	ts := NewTokenService(
		"SUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYAT",
		15*time.Minute,
		30*time.Hour,
	)
	return &Service{
		store:        store,
		tokenService: ts,
	}
}

func (s *Service) Authorize(ctx context.Context, r *api.AuthorizeRequest) (*api.Authorization, error) {
	ju, err := s.tokenService.ParseToken(r.Token)
	if err != nil {
		return nil, err
	}
	return &api.Authorization{
		Id:       ju.ID,
		Username: ju.Username,
	}, nil
}
