package mocks

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type DatabaseMock struct {
	mock.Mock
}

func NewDatabaseMock() *DatabaseMock {
	return &DatabaseMock{}
}

func (d *DatabaseMock) GetDB() *gorm.DB {
	args := d.Called()
	return args.Get(0).(*gorm.DB)
}

func (d *DatabaseMock) WithTransaction(fn func(database.Database) error) error {
	args := d.Called(fn)
	return args.Error(0)
}
