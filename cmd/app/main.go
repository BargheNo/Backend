package main

import (
	"fmt"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/presentation/routes"
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	ginEngine := gin.New()

	config := bootstrap.Run()

	app, err := wire.InitializeApplication(config)
	if err != nil {
		panic(err)
	}
	app.Database.DB.GetDB().AutoMigrate(
		&entity.Address{},
		&entity.Permission{},
		&entity.Role{},
		&entity.User{},
		&entity.Corporation{},
		&entity.ContactInformation{},
		&entity.InstallationRequest{},
		&entity.Bid{},
	)

	routes.Run(ginEngine, app)

	ginEngine.Run(fmt.Sprintf(":%v", config.Env.Server.Port))
}
