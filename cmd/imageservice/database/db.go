package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"image-service/cmd/imageservice/cnf"
	"image-service/cmd/imageservice/database/entities"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

const TmpCollection = "TmpImages"
const PermanentCollection = "PermanentImages"

type Db interface {
	StoreTmpImgInfo(ctx context.Context, tmpImg entities.TmpImage) (string, error)

	StorePermanentImgInfo(ctx context.Context, image entities.PermanentImage) error

	GetTmpImgInfo(ctx context.Context, imageId string) (entities.TmpImage, error)

	CleanExpired(ctx context.Context, notConfirmedAfter time.Duration) (int64, error)
}

type MongoDB struct {
	inner mongo.Database
	Db
}

func NewMongoDB(config cnf.ImageServiceConfig) MongoDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		panic(err)
	}
	return MongoDB{
		inner: *client.Database(config.MongoDatabase),
	}
}

func (db MongoDB) StoreTmpImgInfo(ctx context.Context, tmpImg entities.TmpImage) (string, error) {
	result, err := db.inner.Collection(TmpCollection).InsertOne(ctx, tmpImg)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db MongoDB) StorePermanentImgInfo(ctx context.Context, image entities.PermanentImage) error {
	_, err := db.inner.Collection(PermanentCollection).InsertOne(ctx, image)
	return err
}

func (db MongoDB) GetTmpImgInfo(ctx context.Context, imageId string) (entities.TmpImage, error) {
	oImageId, err := primitive.ObjectIDFromHex(imageId)
	if err != nil {
		return entities.TmpImage{}, err
	}

	result := db.inner.Collection(TmpCollection).FindOne(ctx, oImageId)
	err = result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entities.TmpImage{}, errors.New("該圖片不存在")
		}
		return entities.TmpImage{}, err
	}
	var image = entities.TmpImage{}
	err = result.Decode(&image)
	return image, err
}

func (db MongoDB) CleanExpired(ctx context.Context, notConfirmedAfter time.Duration) (int64, error) {
	var filter = bson.D{}
	result, err := db.inner.Collection(TmpCollection).DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
