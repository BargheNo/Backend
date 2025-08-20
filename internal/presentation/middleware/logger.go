package middleware

import (
	"time"

	"github.com/BargheNo/Backend/internal/domain/logger"
	"github.com/gin-gonic/gin"
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

	if path == "/metrics" {
		c.Next()
		return
	}

	c.Next()

	latency := time.Since(start)

	if len(c.Errors) > 0 {
		for _, e := range c.Errors.Errors() {
			log.logger.Error(
				"Request Error",
				logger.Int("status", c.Writer.Status()),
				logger.String("method", c.Request.Method),
				logger.String("path", path),
				logger.String("query", query),
				logger.String("ip", c.ClientIP()),
				logger.Duration("latency", latency),
				logger.String("user-agent", c.Request.UserAgent()),
				logger.String("error", e),
			)
		}
	} else {
		log.logger.Info(
			"Request",
			logger.Int("status", c.Writer.Status()),
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", query),
			logger.String("ip", c.ClientIP()),
			logger.Duration("latency", latency),
			logger.String("user-agent", c.Request.UserAgent()),
		)
	}
}
