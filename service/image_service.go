package service

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/e-commerce-microservices/image-service/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ImageService ...
type ImageService struct {
	imageStore ImageStore
	pb.UnimplementedImageServiceServer
}

// NewImageService create ImageService instance
func NewImageService() ImageService {
	return ImageService{
		imageStore: NewDiskImageStore("images"),
	}
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
		log.Print("waiting to receive more data")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive chunk data: %v ", err)
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size: %d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize)
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
		}
	}

	imageID, err := srv.imageStore.Save(id, imageType, imageData)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}

	log.Printf("saved image with id: %s, size: %d", id, imageSize)

	return nil
}
