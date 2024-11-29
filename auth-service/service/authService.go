package service

import (
	"awesoma31/common/api"
	"awesoma31/common/storage/model"
	"context"
	"errors"
	"fmt"
	"github.com/awesoma31/auth-service/storage"
	"log"
	"time"
)

type AuthService interface {
	Authorize(ctx context.Context, r *api.AuthorizeRequest) (*api.Authorization, error)
	Login(ctx context.Context, r *api.LoginRequest) (*api.Tokens, error)
	Register(ctx context.Context, r *api.LoginRequest) error
}

type service struct {
	store        storage.UserStore
	tokenService TokenService
}

func NewAuthService(store storage.UserStore) AuthService {
	ts := NewTokenService(
		"SUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYATSUKABLYAT",
		1500*time.Minute,
		30*time.Hour,
	)

	return &service{
		store:        store,
		tokenService: ts,
	}
}

func (s *service) Authorize(ctx context.Context, r *api.AuthorizeRequest) (*api.Authorization, error) {
	ju, err := s.tokenService.ParseToken(ctx, r.Token)
	if err != nil {
		return nil, err
	}
	return &api.Authorization{
		Id:       ju.ID,
		Username: ju.Username,
	}, nil
}

func (s *service) Login(ctx context.Context, r *api.LoginRequest) (*api.Tokens, error) {
	userId, err := s.store.FindIdByUsername(ctx, r.GetUsername())
	if err != nil {
		log.Printf("user not found: %s\n", err)
		return nil, err
	}

	tokens, err := s.tokenService.GenerateTokens(userId, r.GetUsername())
	if err != nil {
		log.Printf("tokens generation failed: %s\n", err)
	}

	return &tokens, nil
}

func (s *service) Register(ctx context.Context, req *api.LoginRequest) error {
	existingUser, err := s.store.FindByUsername(ctx, req.Username)
	if existingUser != nil {
		return errors.New("username already taken")
	}

	if err != nil && !errors.Is(err, storage.ErrUserNotFound) {
		log.Printf("Unexpected error during user lookup: %v", err)
		return fmt.Errorf("internal server error: %w", err)
	}

	user := model.NewUser(req.GetUsername(), req.GetPassword())

	if err := s.store.Create(ctx, user); err != nil {
		return err
	}

	return nil
}
