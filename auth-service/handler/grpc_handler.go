package handler

import (
	pb "awesoma31/common/api"
	"context"
	"github.com/awesoma31/auth-service/service"
	"google.golang.org/grpc"
)

type grpcAuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewGRPCAuthHandler(grpcServer *grpc.Server, svc service.AuthService) {
	handler := &grpcAuthHandler{
		authService: svc,
	}
	pb.RegisterAuthServiceServer(grpcServer, handler)
}

func (h *grpcAuthHandler) Authorize(ctx context.Context, r *pb.AuthorizeRequest) (*pb.Authorization, error) {
	authorization, err := h.authService.Authorize(ctx, r)
	if err != nil {
		return nil, err
	}
	return authorization, nil
}
