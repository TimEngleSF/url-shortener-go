package S3

import (
	"bytes"
	"context"
	"log"

	"github.com/TimEngleSF/url-shortener-go/internal/helper"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Interface interface {
	UploadFile(ctx context.Context, key string, img []byte) (string, error)
}

type S3 struct {
	Client *s3.Client
}

var (
	s3Client *S3
	bucket   string
	region   string
)

func NewS3Client() (*s3.Client, error) {
	// Create an AWS config
	if region == "" || bucket == "" {
		region = helper.GetEnv("AWS_REGION")
		bucket = helper.GetEnv("S3_BUCKET")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	return client, nil
}

func (s *S3) UploadFile(ctx context.Context, key string, img []byte) (string, error) {
	r := bytes.NewReader(img)
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   r,
	})
	if err != nil {
		return "", err
	}

	objectURL := "https://" + bucket + ".s3-" + region + ".amazonaws.com/" + key

	return objectURL, nil
}
