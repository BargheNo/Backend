package mocks

import (
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
