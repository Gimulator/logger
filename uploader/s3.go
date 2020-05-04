package uploader

import (
	"fmt"
	"os"

	minio "github.com/minio/minio-go/v6"
)

var (
	s3URL       string
	s3AccessKey string
	s3SecretKey string
	s3Bucket    string
	s3Key       string
)

type S3 struct {
	*minio.Client
}

func (s *S3) env() error {
	s3URL = os.Getenv("LOGGER_S3_URL")
	if s3URL == "" {
		return fmt.Errorf("set the 'LOGGER_S3_URL' environment variable to connect to s3")
	}
	s3AccessKey = os.Getenv("LOGGER_S3_ACCESS_KEY")
	if s3AccessKey == "" {
		return fmt.Errorf("set the 'S3_ACCESS_KEY' environment variable to connect to s3")
	}
	s3SecretKey = os.Getenv("LOGGER_S3_SECRET_KEY")
	if s3SecretKey == "" {
		return fmt.Errorf("set the 'LOGGER_S3_SECRET_KEY' environment variable to connect to s3")
	}
	s3Bucket = os.Getenv("LOGGER_S3_BUCKET")
	if s3Bucket == "" {
		return fmt.Errorf("set the 'LOGGER_S3_BUCKET' environment variable to upload log file to S3")
	}
	s3Key = os.Getenv("LOGGER_S3_KEY")
	if s3Key == "" {
		return fmt.Errorf("set the 'LOGGER_S3_KEY' environment variable to upload log file to S3")
	}
	return nil
}

func NewS3() (*S3, error) {
	s := &S3{}
	if err := s.env(); err != nil {
		return nil, err
	}

	cli, err := minio.New(s3URL, s3AccessKey, s3SecretKey, false)
	if err != nil {
		return nil, err
	}
	s.Client = cli

	return s, nil
}

func (s *S3) Upload(path string) error {
	_, err := s.FPutObject(s3Bucket, s3Key, path, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
