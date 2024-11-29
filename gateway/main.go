package main

import (
	"awesoma31/common"
	"awesoma31/common/api"
	"awesoma31/gateway/handlers"
	"awesoma31/gateway/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type application struct {
	httpAddr          string
	authServiceAddr   string
	pointsServiceAddr string
}

func main() {
	app := &application{
		httpAddr:          common.GetEnv("HTTP_ADDR", ":8083"),
		authServiceAddr:   common.GetEnv("AUTH_SERVICE_ADDR", "localhost:8081"),
		pointsServiceAddr: common.GetEnv("POINTS_SERVICE_ADDR", "localhost:8082"),
	}

	authConn, err := grpc.NewClient(app.authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer authConn.Close()

	pointsConn, err := grpc.NewClient(app.pointsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	authSvc := api.NewAuthServiceClient(authConn)
	pointsSvc := api.NewPointsServiceClient(pointsConn)
	log.Printf("Dialing auth-service  at %s", app.authServiceAddr)
	log.Printf("Dialing points-service  at %s", app.pointsServiceAddr)

	mux := http.NewServeMux()

	handler := handlers.NewHandler(authSvc, pointsSvc)
	handler.MountRoutes(mux)

	corsMux := middleware.CorsMiddleware(mux)

	log.Printf("Listening on %s\n", app.httpAddr)

	if err := http.ListenAndServe(app.httpAddr, corsMux); err != nil {
		log.Fatal(err)
	}
}
