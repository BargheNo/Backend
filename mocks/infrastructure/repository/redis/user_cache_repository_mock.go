package mocks

import (
	"time"

	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type UserCacheRepositoryMock struct {
	mock.Mock
}

func NewUserCacheRepositoryMock() *UserCacheRepositoryMock {
	return &UserCacheRepositoryMock{}
}

func (u *UserCacheRepositoryMock) SetUser(user *entity.User, expiration time.Duration) error {
	args := u.Called(user, expiration)
	return args.Error(0)
}

func (u *UserCacheRepositoryMock) GetUser(userID uint) (*entity.User, error) {
	args := u.Called(userID)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserCacheRepositoryMock) DeleteUser(userID uint) error {
	args := u.Called(userID)
	return args.Error(0)
}

func (u *UserCacheRepositoryMock) SetUserToken(userID uint, token string, expiration time.Duration) error {
	args := u.Called(userID, token, expiration)
	return args.Error(0)
}

func (u *UserCacheRepositoryMock) GetUserToken(userID uint) (string, error) {
	args := u.Called(userID)
	return args.String(0), args.Error(1)
}

func (u *UserCacheRepositoryMock) DeleteUserToken(userID uint) error {
	args := u.Called(userID)
	return args.Error(0)
}
