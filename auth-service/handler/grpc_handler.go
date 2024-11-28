package handler

import (
	"awesoma31/common/api"
	"context"
	"github.com/awesoma31/auth-service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcAuthHandler struct {
	api.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewGRPCAuthHandler(grpcServer *grpc.Server, svc service.AuthService) {
	handler := &grpcAuthHandler{
		authService: svc,
	}
	api.RegisterAuthServiceServer(grpcServer, handler)
}

func (h *grpcAuthHandler) Authorize(ctx context.Context, r *api.AuthorizeRequest) (*api.Authorization, error) {
	authorization, err := h.authService.Authorize(ctx, r)
	if err != nil {
		return nil, err
	}
	return authorization, nil
}

func (h *grpcAuthHandler) Login(ctx context.Context, r *api.LoginRequest) (*api.Tokens, error) {
	tokens, err := h.authService.Login(ctx, r)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (h *grpcAuthHandler) Register(ctx context.Context, r *api.LoginRequest) (*emptypb.Empty, error) {
	if r.GetUsername() == "" || r.GetPassword() == "" {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "username and password cannot be empty")
	}

	err := h.authService.Register(ctx, r)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
