package main

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"image-service/cmd/imageservice/database"
	"image-service/cmd/imageservice/database/entities"
	imageService "image-service/protobuffs"
	"time"
)

type GrpcImageService struct {
	database database.Db
	imageService.ImageServiceServer
}

func (s GrpcImageService) ConfirmImageUse(ctx context.Context, req *imageService.ConfirmImageUsedRequest) (*imageService.ConfirmImageUsedResponse, error) {
	tmpImage, err := s.database.GetTmpImgInfo(ctx, req.ImageId)

	if err == database.ErrNotFound {
		return nil, status.Errorf(codes.NotFound, "Could not find image %s", req.ImageId)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req.Strict && (*req.Usage != tmpImage.ExpectedUsage || *req.UserId != tmpImage.Uploader.Hex()) {
		return nil, status.Error(codes.Unauthenticated, "wrong usage or the image doesn't belong to you")
	}

	err = s.database.StorePermanentImgInfo(ctx, entities.NewPermanentImage(tmpImage))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &imageService.ConfirmImageUsedResponse{
		Success: true,
	}, nil
}

func (s GrpcImageService) CleanExpiredTmp(ctx context.Context, _ *imageService.Empty) (*imageService.CleanExpiredTmpResponse, error) {
	cleaned, err := s.database.CleanExpired(ctx, time.Duration(time.Hour*24*30))
	return &imageService.CleanExpiredTmpResponse{
		Cleaned: cleaned,
	}, err
}
