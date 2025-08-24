package mocks

import (
	"context"
	"time"

	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/stretchr/testify/mock"
)

type UserCacheRepositoryMock struct {
	mock.Mock
}

func NewUserCacheRepositoryMock() *UserCacheRepositoryMock {
	return &UserCacheRepositoryMock{}
}

func (u *UserCacheRepositoryMock) Get(ctx context.Context, key string) (*userdto.OTPData, error) {
	args := u.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userdto.OTPData), args.Error(1)
}

func (u *UserCacheRepositoryMock) Set(ctx context.Context, key, otp string, expiration time.Duration) error {
	args := u.Called(ctx, key, otp, expiration)
	return args.Error(0)
}
