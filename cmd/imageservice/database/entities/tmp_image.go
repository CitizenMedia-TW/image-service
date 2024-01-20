package entities

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"image-service/protobuffs/image-service"
	"time"
)

type TmpImage struct {
	Id            primitive.ObjectID       `bson:"_id"`
	Path          string                   `bson:"path"`
	Uploader      primitive.ObjectID       `bson:"uploader"`
	UploadedAt    primitive.DateTime       `bson:"uploaded_at"`
	ExpectedUsage image_service.ImageUsage `bson:"expected_usage"`
	Used          bool                     `bson:"used"`
}

func NewTmpImage(uploader string, usage image_service.ImageUsage, path string) (TmpImage, error) {

	uOid, err := primitive.ObjectIDFromHex(uploader)

	if err != nil {
		return TmpImage{}, errors.Join(errors.New("uploader id 非正確格式: "), err)
	}

	return TmpImage{
		Id:            primitive.NewObjectID(),
		Path:          path,
		UploadedAt:    primitive.NewDateTimeFromTime(time.Now()),
		ExpectedUsage: usage,
		Uploader:      uOid,
	}, nil
}
