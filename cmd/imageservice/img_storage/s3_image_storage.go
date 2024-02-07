package img_storage

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"image-service/cmd/imageservice/cnf"
	"time"
)

type S3ImageStorage struct {
	svc    *s3.S3
	config cnf.ImageServiceConfig
	ImageStorage
}

func (m S3ImageStorage) Store(fileSuffix string, imageData []byte, imageId string) (string, error) {
	now := time.Now()
	if imageId == "" {
		uid, err := uuid.NewRandom()
		if err != nil {
			return "", nil
		}
		imageId = uid.String()
	}
	key := fmt.Sprintf("/images/%d/%d/%d/%s.%s", now.Year(), now.Month(), now.Day(), imageId, fileSuffix)
	_, err := m.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(m.config.S3Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(imageData),
	})
	if err != nil {
		return "", errors.Join(errors.New("could not upload to s3"), err)
	}
	return key, nil
}

func (m S3ImageStorage) Delete(path string) error {
	_, err := m.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(m.config.S3Bucket),
		Key:    aws.String(path),
	})
	return err
}

func NewS3ImageStorage(config cnf.ImageServiceConfig) S3ImageStorage {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:                   credentials.NewStaticCredentials(config.S3KeyId, config.S3KeyValue, ""),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Region:                        aws.String(config.AwsRegion),
	}))
	svc := s3.New(sess)
	return S3ImageStorage{
		svc:    svc,
		config: config,
	}
}
