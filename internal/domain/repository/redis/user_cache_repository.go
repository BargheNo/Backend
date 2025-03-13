package cacherepository

import (
	"context"
	"time"

	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type UserCacheRepository interface {
	Get(ctx context.Context, key string) (*userdto.OTPData, bool)
	Set(ctx context.Context, key, otp string, expiration time.Duration) error
}
