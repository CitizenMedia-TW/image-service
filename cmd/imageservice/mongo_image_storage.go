package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type MongoImageStorage struct {
	database *mongo.Database
	config   ImageServiceConfig
}

// MongoFields is the struct that defines the fields in the MongoDB database
type MongoFields struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Data     []byte             `bson:"data" json:"data"`
	Type     string             `bson:"type" json:"type"`
	Uploaded time.Time          `bson:"uploaded" json:"uploaded"`
}

func (m *MongoImageStorage) Store(fileName string, imageData []byte) (string, error) {
	imageType := http.DetectContentType(imageData)

	image := MongoFields{
		Name:     fileName,
		Data:     imageData,
		Type:     imageType,
		Uploaded: time.Now(),
	}

	// Insert the image into the database
	storedImage, err := m.database.Collection(m.config.mongoCollection).InsertOne(context.Background(), image)
	if err != nil {
		return "", errors.New("error inserting image data into MongoDB")
	}
	return m.config.host + "/" + storedImage.InsertedID.(primitive.ObjectID).Hex(), nil
}

func NewMongoImageStorage(config ImageServiceConfig) ImageStorage {
	// connectToDB connects to the MongoDB database and returns the collection
	log.Println("Connecting to MongoDB...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use this for when running in Docker
	// clientOptions := options.Client().ApplyURI("mongodb://root:rootpassword@mongo:27017/")
	// client, err := mongo.Connect(ctx, clientOptions)

	// Use this for when running locally
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.mongoURI))

	// Connect to MongoDB Atlas
	// mongodbURI := os.Getenv("MONGODB_URI")
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURI))

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	// collection := client.Database("GolangImageTest").Collection("images")

	storage := MongoImageStorage{
		database: client.Database(config.mongoDatabase),
	}

	return &storage

}
