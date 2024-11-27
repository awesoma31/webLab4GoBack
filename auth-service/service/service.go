package service

import (
	"github.com/awesoma31/auth-service/storage"
)

type Service struct {
	store storage.UsersStore
}

func NewAuthService(store storage.UsersStore) *Service {
	return &Service{store: store}
}
