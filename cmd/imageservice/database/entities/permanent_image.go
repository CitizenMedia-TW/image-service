package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	imageService "image-service/protobuffs/image-service"
)

type PermanentImage struct {
	Id         primitive.ObjectID      `bson:"_id"`
	Path       string                  `bson:"path"`
	Uploader   primitive.ObjectID      `bson:"uploader"`
	UploadedAt primitive.DateTime      `bson:"uploaded_at"`
	Usage      imageService.ImageUsage `bson:"expected_usage"`
}

func NewPermanentImage(image TmpImage) PermanentImage {
	return PermanentImage{
		Id:         image.Id,
		Path:       image.Path,
		Uploader:   image.Uploader,
		UploadedAt: image.UploadedAt,
		Usage:      image.ExpectedUsage,
	}
}
