package api

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"image-service/cmd/imageservice/database"
	"image-service/cmd/imageservice/database/entities"
	"image-service/cmd/imageservice/img_storage"
	imageService "image-service/protobuffs/image-service"
	"log"
	"net"
	"time"
)

type GrpcImageService struct {
	database     database.Db
	imageStorage img_storage.ImageStorage
	imageService.ImageServiceServer
}

func StartGrpc(storage img_storage.ImageStorage, db database.Db) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 1111))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	imageService.RegisterImageServiceServer(grpcServer, GrpcImageService{
		imageStorage: storage,
		database:     db,
	})

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func (s GrpcImageService) ConfirmImageUse(ctx context.Context, req *imageService.ConfirmImageUsedRequest) (*imageService.ConfirmImageUsedResponse, error) {
	tmpImage, err := s.database.GetTmpImgInfo(ctx, req.ImageId)

	if err == database.ErrNotFound {
		return nil, status.Errorf(codes.NotFound, "Could not find image %s", req.ImageId)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "error getting tmp img from database: "+err.Error())
	}

	if req.Strict && (*req.Usage != tmpImage.ExpectedUsage || *req.UserId != tmpImage.Uploader.Hex()) {
		return nil, status.Error(codes.Unauthenticated, "wrong usage or the image doesn't belong to you")
	}

	err = s.database.StorePermanentImgInfo(ctx, entities.NewPermanentImage(tmpImage))
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Join(errors.New("error while storing permanent image"), err).Error())
	}
	return &imageService.ConfirmImageUsedResponse{
		Success: true,
	}, nil
}

func (s GrpcImageService) CleanExpiredTmp(ctx context.Context, _ *imageService.Empty) (*imageService.CleanExpiredTmpResponse, error) {
	cleaned, err := s.database.CleanExpired(ctx, time.Hour*24*30)
	return &imageService.CleanExpiredTmpResponse{
		Cleaned: cleaned,
	}, err
}

func (s GrpcImageService) DeleteImage(ctx context.Context, r *imageService.DeleteImageRequest) (*imageService.DeleteImageResponse, error) {
	deletedPath, err := s.database.DeletePermanent(ctx, r.ImageId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if deletedPath == "" {
		return nil, status.Error(codes.NotFound, "Could not find image with correspond Id")
	}

	err = s.imageStorage.Delete(deletedPath)
	if err != nil {
		return nil, err
	}

	return &imageService.DeleteImageResponse{
		Success: err == nil,
	}, err
}
