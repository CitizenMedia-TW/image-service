package cnf

import "os"

type ImageServiceConfig struct {
	StorageType string
	Host        string

	S3Bucket   string
	S3KeyId    string
	S3KeyValue string
	AwsRegion  string

	MongoURI      string
	MongoDatabase string
	host          string
}

func NewConfig() ImageServiceConfig {
	host := os.Getenv("HOST")
	s3Bucket := os.Getenv("S3_BUCKET")
	s3KeyValue := os.Getenv("S3_KEY_VALUE")
	s3KeyId := os.Getenv("S3_KEY_ID")
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDatabase := os.Getenv("MONGODB_DBNAME")
	awsRegion := os.Getenv("AWS_REGION")
	return ImageServiceConfig{
		Host:          host,
		S3Bucket:      s3Bucket,
		S3KeyValue:    s3KeyValue,
		S3KeyId:       s3KeyId,
		MongoURI:      mongoURI,
		MongoDatabase: mongoDatabase,
		AwsRegion:     awsRegion,
	}
}
