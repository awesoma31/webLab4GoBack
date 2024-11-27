package main

import (
	"awesoma31/common"
	"github.com/awesoma31/auth-service/handler"
	"github.com/awesoma31/auth-service/service"
	"github.com/awesoma31/auth-service/storage"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	grpcAddr = common.GetEnv("GRPC_ADDR", "localhost:8081")
)

func main() {
	store := storage.NewStore()
	svc := service.NewAuthService(store)

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {

		log.Fatal(err)
	}
	defer l.Close()

	handler.NewGRPCAuthHandler(grpcServer, svc)

	log.Printf("Listening GRPC on: %s", grpcAddr)

	if err = grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
