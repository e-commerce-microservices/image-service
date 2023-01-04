package main

import (
	"log"
	"net"

	"github.com/e-commerce-microservices/image-service/pb"
	"github.com/e-commerce-microservices/image-service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	grpcServer := grpc.NewServer()

	imageService := service.NewImageService()
	pb.RegisterImageServiceServer(grpcServer, imageService)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	reflection.Register(grpcServer)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
