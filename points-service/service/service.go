package service

import (
	"github.com/awesoma31/points-service/storage"
)

type PointsService struct {
	store storage.PointsStore
}

func NewPointsService(store storage.PointsStore) *PointsService {
	return &PointsService{store: store}
}
