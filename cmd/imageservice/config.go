package main

import "os"

type ImageServiceConfig struct {
	storageType string
	host        string

	s3Bucket  string
	s3Key     string
	awsRegion string

	mongoURI        string
	mongoDatabase   string
	mongoCollection string
}

func NewConfig() ImageServiceConfig {
	storageType := os.Getenv("STORAGE_TYPE")
	host := os.Getenv("HOST")
	s3Bucket := os.Getenv("S3_BUCKET")
	s3Key := os.Getenv("S3_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDatabase := os.Getenv("MONGODB_DATABASE")
	mongoCollection := os.Getenv("MONGODB_COLLECTION")
	return ImageServiceConfig{
		storageType:     storageType,
		host:            host,
		s3Bucket:        s3Bucket,
		s3Key:           s3Key,
		awsRegion:       awsRegion,
		mongoURI:        mongoURI,
		mongoDatabase:   mongoDatabase,
		mongoCollection: mongoCollection,
	}
}
