package database

import (
	"fmt"
	"sync"

	"github.com/BargheNo/Backend/bootstrap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	GetDB() *gorm.DB
}

type PostgresDatabase struct {
	DB *gorm.DB
}

var (
	dbOnce     sync.Once
	dbInstance *PostgresDatabase
)

func NewPostgresDatabase(dbConfig *bootstrap.Database) *PostgresDatabase {
	dbOnce.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			dbConfig.Host,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
			dbConfig.Port,
		)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("failed to connect database"))
		}

		dbInstance = &PostgresDatabase{DB: db}
	})

	return dbInstance
}

func (pgx *PostgresDatabase) GetDB() *gorm.DB {
	return dbInstance.DB
}
