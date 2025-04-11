package main

import (
	"fmt"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
	"github.com/BargheNo/Backend/internal/presentation/routes"
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	ginEngine := gin.New()

	config := bootstrap.Run()

	hub := websocket.NewHub()
	go hub.Run()

	app, err := wire.InitializeApplication(config, hub)
	if err != nil {
		panic(err)
	}

	app.Database.DB.GetDB().AutoMigrate(
		&entity.Address{},
		&entity.City{},
		&entity.ChatMessage{},
		&entity.ChatRoom{},
		&entity.NotificationSetting{},
		&entity.NotificationType{},
		&entity.Notification{},
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
		&entity.Panel{},
		&entity.MaintenanceRequest{},
		&entity.MaintenanceRecord{},
	)

	app.Seeds.AddressSeeder.SeedProvincesAndCities()
	app.Seeds.NotificationTypeSeeder.SeedNotificationTypes()

	routes.Run(ginEngine, app)

	ginEngine.Run(fmt.Sprintf(":%v", config.Env.Server.Port))
}
