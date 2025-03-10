package database

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	GetDB() *gorm.DB
}

type postgresDatabase struct {
	DB *gorm.DB
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

func NewPostgresDatabase(dsn string) Database {
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("failed to connect database"))
		}

		dbInstance = &postgresDatabase{DB: db}
	})

	return dbInstance
}

func (pgx *postgresDatabase) GetDB() *gorm.DB {
	return dbInstance.DB
}
