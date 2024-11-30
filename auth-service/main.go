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
	host := common.GetEnv("DB_HOST", "localhost")
	port := common.GetEnv("DB_PORT", "5432")
	user := common.GetEnv("DB_USER", "awesoma")
	password := common.GetEnv("DB_PASSWORD", "1")
	dbname := common.GetEnv("DB_NAME", "lab4")

	store := storage.NewUserStore(
		storage.WithHost(host),
		storage.WithPort(port),
		storage.WithUsername(user),
		storage.WithPassword(password),
		storage.WithDBName(dbname),
	)
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
