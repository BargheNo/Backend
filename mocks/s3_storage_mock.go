package mocks

import (
	"mime/multipart"
	"time"

	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/stretchr/testify/mock"
)

type S3StorageMock struct {
	mock.Mock
}

func NewS3StorageMock() *S3StorageMock {
	return &S3StorageMock{}
}

func (s *S3StorageMock) DeleteObject(bucketType enum.BucketType, key string) error {
	args := s.Called(bucketType, key)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (s *S3StorageMock) GetPresignedURL(bucketType enum.BucketType, objectKey string, expiration time.Duration) string {
	args := s.Called(bucketType, objectKey, expiration)
	if args.Get(0) != nil {
		return args.Get(0).(string)
	}
	return ""
}

func (s *S3StorageMock) UploadObject(bucketType enum.BucketType, key string, file *multipart.FileHeader) {
	s.Called(bucketType, key, file)
}
