package middleware

import (
	"time"

	"github.com/BargheNo/Backend/internal/domain/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoggerMiddleware struct {
	logger logger.Logger
}

func NewLoggerMiddleware(logger logger.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
	}
}

func (log *LoggerMiddleware) GinLoggerMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	c.Next()

	latency := time.Since(start)

	if len(c.Errors) > 0 {
		for _, e := range c.Errors.Errors() {
			log.logger.Error(
				"Request Error",
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.Duration("latency", latency),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("error", e),
			)
		}
	} else {
		log.logger.Info(
			"Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	}
}
