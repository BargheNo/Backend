package mocks

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type DatabaseMock struct {
	mock.Mock
	db *gorm.DB
}

func NewDatabaseMock() *DatabaseMock {
	return &DatabaseMock{}
}

func (m *DatabaseMock) GetDB() *gorm.DB {
	return m.db
}

func (m *DatabaseMock) WithTransaction(fn func(tx database.Database) error) error {
	args := m.Called(fn)
	if fn != nil {
		_ = fn(m)
	}
	return args.Error(0)
}
