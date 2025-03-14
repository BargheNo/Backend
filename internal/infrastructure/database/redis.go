package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetRDB() *redis.Client
}

type RedisDatabase struct {
	RDB *redis.Client
}

var (
	rdbOnce     sync.Once
	rdbInstance *RedisDatabase
)

func NewRedisDatabase(redisConfig *bootstrap.Redis) *RedisDatabase {
	rdbOnce.Do(func() {
		rdbNumber, _ := strconv.Atoi(redisConfig.RDBNumber)
		address := fmt.Sprintf("%s:%s", redisConfig.Address, redisConfig.Port)
		rdb := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: redisConfig.Password,
			DB:       rdbNumber,
		})
		_, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal("Error connecting to Redis:", err)
		}
		rdbInstance = &RedisDatabase{RDB: rdb}
	})

	return rdbInstance
}

func (rdb *RedisDatabase) GetRDB() *redis.Client {
	return rdbInstance.RDB
}
