package handler

import (
	pb "awesoma31/common/api"
	"context"
	"github.com/awesoma31/points-service/service"
	"google.golang.org/grpc"
	"log"
)

type grpcPointsHandler struct {
	pb.UnimplementedPointsServiceServer
	pointsService *service.PointsService
}

func NewGRPCPointsHandler(grpcServer *grpc.Server, ps *service.PointsService) {
	handler := &grpcPointsHandler{pointsService: ps}
	pb.RegisterPointsServiceServer(grpcServer, handler)
}

func (g *grpcPointsHandler) GetUserPointsPage(ctx context.Context, r *pb.PointsPageRequest) (*pb.PointsPage, error) {
	log.Println("user point page request received")
	pp := &pb.PointsPage{
		Content:       make([]*pb.Point, 0),
		PageNumber:    0,
		PageSize:      0,
		TotalElements: 0,
		TotalPages:    0,
	}
	return pp, nil
}
