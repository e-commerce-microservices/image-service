package main

import (
	"encoding/base64"
	"log"
	"net"
	"net/http"

	"github.com/e-commerce-microservices/image-service/pb"
	"github.com/e-commerce-microservices/image-service/service"
	"google.golang.org/grpc"
)

func main() {
	go serveFileServer()

	grpcServer := grpc.NewServer()

	imageService := service.NewImageService()
	pb.RegisterImageServiceServer(grpcServer, imageService)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}

}

func serveFileServer() {
	fs := http.FileServer(http.Dir("./images"))
	http.Handle("/", fs)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println(err)
	}

}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func toBytes(str string) []byte {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
