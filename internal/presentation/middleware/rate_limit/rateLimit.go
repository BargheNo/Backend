package ratelimit

import (
	"strconv"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimitMiddleware struct {
	rateLimit *bootstrap.RateLimit
}

func NewRateLimit(rateLimit *bootstrap.RateLimit) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		rateLimit: rateLimit,
	}
}

func (rl *RateLimitMiddleware) RateLimit(c *gin.Context) {
	limit, err := strconv.Atoi(rl.rateLimit.Limit)
	if err != nil {
		limit = 5
	}
	burst, err := strconv.Atoi(rl.rateLimit.Burst)
	if err != nil {
		burst = 10
	}
	limiter := rate.NewLimiter(rate.Limit(limit), burst)
	if !limiter.Allow() {
		rateLimitError := exception.NewRateLimitError()
		panic(rateLimitError)
	}
	c.Next()
}
