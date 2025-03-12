package cacherepositoryimpl

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type UserCacheRepository struct {
	rdb database.Cache
}

func NewUserCacheRepository(rdb database.Cache) *UserCacheRepository {
	return &UserCacheRepository{
		rdb: rdb,
	}
}
