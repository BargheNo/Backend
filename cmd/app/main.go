package main

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	ginEngine := gin.Default()

	config := bootstrap.Run()

	ginEngine.Run(config.Env.Server.Port)
}
