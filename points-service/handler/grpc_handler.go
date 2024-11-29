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
	ps service.PointsService
}

func NewGRPCPointsHandler(grpcServer *grpc.Server, ps service.PointsService) {
	handler := &grpcPointsHandler{ps: ps}
	pb.RegisterPointsServiceServer(grpcServer, handler)
}

func (h *grpcPointsHandler) GetUserPointsPage(ctx context.Context, r *pb.PointsPageRequest) (*pb.PointsPage, error) {
	log.Println("user point page request received")

	p, err := h.ps.GetPointsPageByID(ctx, r.GetPageParam(), r.PageSize, r.Id)
	if err != nil {
		return nil, err
	}

	return p, nil
}
