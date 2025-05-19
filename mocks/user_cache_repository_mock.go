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

func (u *UserCacheRepositoryMock) Get(ctx context.Context, key string) (*userdto.OTPData, bool) {
	args := u.Called(ctx, key)
	return args.Get(0).(*userdto.OTPData), args.Bool(1)
}

func (u *UserCacheRepositoryMock) Set(ctx context.Context, key, otp string, expiration time.Duration) error {
	args := u.Called(ctx, key, otp, expiration)
	return args.Error(0)
}
