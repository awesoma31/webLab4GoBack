package storage

import "context"

type PointsStore interface {
	GetUserPointsPage(ctx context.Context) error
}

type Store struct {
	//todo
}

func NewStore() *Store {
	//todo psql
	return &Store{}
}

func (s *Store) GetUserPointsPage(ctx context.Context) error {
	//TODO implement me
	return nil
}
