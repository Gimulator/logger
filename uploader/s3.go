package uploader

import (
	"fmt"
	"os"

	minio "github.com/minio/minio-go/v6"
)

type S3 struct {
	*minio.Client

	s3Url       string
	s3AccessKey string
	s3SecretKey string
	s3Bucket    string
}

func NewS3() (*S3, error) {
	s := &S3{}
	if err := s.env(); err != nil {
		return nil, err
	}

	cli, err := minio.New(s.s3Url, s.s3AccessKey, s.s3SecretKey, false)
	if err != nil {
		return nil, err
	}
	s.Client = cli

	return s, nil
}

func (s *S3) Upload(path, key string) error {
	_, err := s.FPutObject(s.s3Bucket, key, path, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *S3) env() error {
	s.s3Url = os.Getenv("LOGGER_S3_URL")
	if s.s3Url == "" {
		return fmt.Errorf("set the 'LOGGER_S3_URL' environment variable to connect to s3")
	}
	s.s3AccessKey = os.Getenv("LOGGER_S3_ACCESS_KEY")
	if s.s3AccessKey == "" {
		return fmt.Errorf("set the 'LOGGER_S3_ACCESS_KEY' environment variable to connect to s3")
	}
	s.s3SecretKey = os.Getenv("LOGGER_S3_SECRET_KEY")
	if s.s3SecretKey == "" {
		return fmt.Errorf("set the 'LOGGER_S3_SECRET_KEY' environment variable to connect to s3")
	}
	s.s3Bucket = os.Getenv("LOGGER_S3_BUCKET")
	if s.s3Bucket == "" {
		return fmt.Errorf("set the 'LOGGER_S3_BUCKET' environment variable to upload log file to S3")
	}
	return nil
}
