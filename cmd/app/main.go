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
		&entity.City{},
		&entity.Permission{},
		&entity.Province{},
		&entity.Role{},
		&entity.CorporationStaff{},
		&entity.Bid{},
		&entity.User{},
		&entity.InstallationRequest{},
		&entity.Corporation{},
		&entity.Signatory{},
		&entity.ContactType{},
		&entity.ContactInformation{},
	)

	app.Seeds.AddressSeeder.SeedProvincesAndCities()

	routes.Run(ginEngine, app)

	ginEngine.Run(fmt.Sprintf(":%v", config.Env.Server.Port))
}
