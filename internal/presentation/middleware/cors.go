package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CORSMiddleware struct{}

func NewCorsMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

func (cm *CORSMiddleware) CORS() gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://185.110.189.68:3001", "http://46.249.99.69:3001"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(corsConfig)
}
