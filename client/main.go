package main

import (
	"context"
	"encoding/base64"
	"log"
	"strings"

	"github.com/e-commerce-microservices/image-service/pb"
	"google.golang.org/grpc"
)

func main() {
	// dial
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewImageServiceClient(conn)
	stream, err := client.UploadImage(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	imageStr := ""
	// send mime type
	tmp := strings.Split(imageStr, "data:image/")
	mimeType := strings.Split(tmp[1], ";")[0]
	err = stream.Send(&pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				ImageType: mimeType,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// send data
	dataChunk := strings.Split(imageStr, ",")[1]
	stream.Send(&pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_ChunkData{
			ChunkData: toBytes(dataChunk),
		},
	})

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.GetImageUrl())

}
func toBytes(str string) []byte {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
