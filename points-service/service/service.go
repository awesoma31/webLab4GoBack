package service

import (
	"context"
	"github.com/awesoma31/points-service/storage"
)

type PointsService struct {
	store storage.PointsStore
}

func NewPointsPointsService(store storage.PointsStore) *PointsService {
	return &PointsService{store: store}
}

func (p *PointsService) GetUserPointsPage(ctx context.Context) error {
	//TODO implement me
	return nil
}
