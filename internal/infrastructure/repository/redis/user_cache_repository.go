package redis

import (
	"context"
	"encoding/json"
	"time"

	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type UserCacheRepository struct {
	rdb database.Cache
}

func NewUserCacheRepository(rdb database.Cache) *UserCacheRepository {
	return &UserCacheRepository{
		rdb: rdb,
	}
}

func (userCache *UserCacheRepository) Get(ctx context.Context, key string) (*userdto.OTPData, error) {
	value, err := userCache.rdb.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var otpData userdto.OTPData
	if err = json.Unmarshal([]byte(value), &otpData); err != nil {
		return nil, err
	}

	return &otpData, nil

}

func (userCache *UserCacheRepository) Set(ctx context.Context, key, otp string, expiration time.Duration) error {
	otpData := userdto.OTPData{
		OTP:      otp,
		Attempts: 0,
	}
	value, err := json.Marshal(otpData)
	if err != nil {
		return err
	}
	err = userCache.rdb.GetRDB().Set(ctx, key, string(value), expiration).Err()
	if err != nil {

		return err
	}
	return nil
}
