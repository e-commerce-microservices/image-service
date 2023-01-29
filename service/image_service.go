package service

import (
	"bytes"
	"context"
	"io"

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
	return &pb.Pong{
		Message: "pong",
	}, nil
}

const maxImageSize = 1 << 20 // 1MB

// UploadImage client streaming
func (srv ImageService) UploadImage(stream pb.ImageService_UploadImageServer) error {

	req, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}

	imageType := req.GetInfo().GetImageType()

	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive chunk data: %v ", err)
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		imageSize += size
		if imageSize > maxImageSize {
			return status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize)
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
		}
	}

	imageURL, err := saveImage(imageType, imageData)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
	}

	res := &pb.UploadImageResponse{
		ImageUrl: imageURL,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}

	return nil
}
