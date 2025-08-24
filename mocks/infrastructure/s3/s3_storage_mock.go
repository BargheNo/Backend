package mocks

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"time"
)

type S3StorageMock struct {
	mock.Mock
}

func NewS3StorageMock() *S3StorageMock {
	return &S3StorageMock{}
}

func (s *S3StorageMock) DeleteObject(bucketType enum.BucketType, key string) error {
	args := s.Called(bucketType, key)
	return args.Error(0)
}

func (s *S3StorageMock) GetPresignedURL(bucketType enum.BucketType, objectKey string, expiration time.Duration) (string, error) {
	args := s.Called(bucketType, objectKey, expiration)
	return args.String(0), args.Error(1)
}

func (s *S3StorageMock) UploadObject(bucketType enum.BucketType, key string, file *multipart.FileHeader) error {
	args := s.Called(bucketType, key, file)
	return args.Error(0)
}
