package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"time"
)

type S3ImageStorage struct {
	svc    *s3.S3
	config ImageServiceConfig
}

func (m S3ImageStorage) Store(fileSuffix string, imageData []byte) (string, error) {
	now := time.Now()
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%d/%d/%d/%s.%s", now.Year(), now.Month(), now.Day(), uid, fileSuffix)
	_, err = m.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(m.config.s3Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(imageData),
	})
	if err != nil {
		return "", err
	}
	return m.config.host + "/images/" + key, nil
}

// NewS3ImageStorage
// aws 好像有很多種驗證credential的方式
// 看你想採用哪一種，再來這邊做調整
func NewS3ImageStorage(config ImageServiceConfig) S3ImageStorage {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	return S3ImageStorage{
		svc: svc,
	}
}
