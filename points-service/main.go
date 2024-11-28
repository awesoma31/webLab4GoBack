package main

import (
	"awesoma31/common"
	"github.com/awesoma31/points-service/handler"
	"github.com/awesoma31/points-service/service"
	"github.com/awesoma31/points-service/storage"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	grpcAddr = common.GetEnv("POINTS_SERVICE_GRPC_ADDR", "localhost:8082")
)

func main() {
	store := storage.NewStore()
	pointsService := service.NewPointsService(store)

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {

		log.Fatal(err)
	}
	defer l.Close()

	handler.NewGRPCPointsHandler(grpcServer, pointsService)
	log.Printf("Listening GRPC on: %s", grpcAddr)

	if err = grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
