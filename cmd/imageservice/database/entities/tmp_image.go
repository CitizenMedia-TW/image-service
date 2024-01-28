package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"image-service/protobuffs/image-service"
	"time"
)

type TmpImage struct {
	Id            primitive.ObjectID       `bson:"_id"`
	Path          string                   `bson:"path"`
	UploadedAt    primitive.DateTime       `bson:"uploaded_at"`
	ExpectedUsage image_service.ImageUsage `bson:"expected_usage"`
	Used          bool                     `bson:"used"`
}

func NewTmpImage(usage image_service.ImageUsage, path string) TmpImage {

	return TmpImage{
		Id:            primitive.NewObjectID(),
		Path:          path,
		UploadedAt:    primitive.NewDateTimeFromTime(time.Now()),
		ExpectedUsage: usage,
	}
}
