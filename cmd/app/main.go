package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	ginEngine := gin.Default()

	ginEngine.Run(":8080")
}
