package handler

import (
	pb "awesoma31/common/api"
	"context"
	"google.golang.org/grpc"
	"log"
)

type grpcAuthHandler struct {
	pb.UnimplementedAuthServiceServer
}

func NewGRPCAuthHandler(grpcServer *grpc.Server) {
	handler := &grpcAuthHandler{}
	pb.RegisterAuthServiceServer(grpcServer, handler)
}

func (h *grpcAuthHandler) Authorize(context.Context, *pb.AuthorizeRequest) (*pb.Authorization, error) {
	log.Println("Authorize request received")
	//todo
	o := &pb.Authorization{
		Id:       1,
		Username: "test_user",
	}
	return o, nil
}
