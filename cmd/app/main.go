package main

import (
	"fmt"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/BargheNo/Backend/internal/presentation/routes"
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	ginEngine := gin.New()

	config := bootstrap.Run()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.Env.PrimaryDB.Host,
		config.Env.PrimaryDB.User,
		config.Env.PrimaryDB.Password,
		config.Env.PrimaryDB.Name,
		config.Env.PrimaryDB.Port,
	)

	db := database.NewPostgresDatabase(dsn)
	db.GetDB().AutoMigrate(
		&entity.Address{},
		&entity.Permission{},
		&entity.Role{},
		&entity.User{},
	)

	app, err := wire.InitializeApplication(config, db)
	if err != nil {
		panic(err)
	}
	routes.Run(ginEngine, app)

	ginEngine.Run(fmt.Sprintf(":%v", config.Env.Server.Port))
}
