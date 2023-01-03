package service

import (
	"bytes"
	"context"
	"log"

	"github.com/e-commerce-microservices/image-service/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ImageService ...
type ImageService struct {
	pb.UnimplementedImageServiceServer
}

// NewImageService create ImageService instance
func NewImageService() ImageService {
	return ImageService{}
}

// Ping pong
func (srv ImageService) Ping(context.Context, *empty.Empty) (*pb.Pong, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}

const maxImageSize = 1 << 20 // 1MB

// CreateImage client streaming
func (srv ImageService) CreateImage(stream pb.ImageService_CreateImageServer) error {

	req, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}

	id := req.GetInfo().GetId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload-image request for laptop %s with image type %s", id, imageType)

	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size: %d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize)
		}

		_, err := imageData.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
		}
	}

	res := &pb.UploadImageResponse{
		Id:   id,
		Size: maxImageSize,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}

	log.Printf("saved image with id: %s, size: %d", id, imageSize)
	return nil
}
