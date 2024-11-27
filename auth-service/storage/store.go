package storage

import "context"

type UsersStore interface {
	GetByUsername(ctx context.Context) error
}

type Store struct {
	//todo
}

func NewStore() *Store {
	//todo psql
	return &Store{}
}

func (s *Store) GetByUsername(ctx context.Context) error {
	return nil
}
