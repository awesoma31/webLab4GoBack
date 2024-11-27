package model

import (
	pb "awesoma31/common/api"
	"context"
)

type PointsService interface {
	GetUserPointsPage(context.Context, *pb.PointsPageRequest) (*pb.PointsPage, error)
}
