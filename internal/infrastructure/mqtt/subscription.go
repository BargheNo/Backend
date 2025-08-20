package mqtt

import (
	"fmt"
	"time"

	"github.com/BargheNo/Backend/internal/application/service"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	database "github.com/BargheNo/Backend/internal/infrastructure/database"
	loggerImpl "github.com/BargheNo/Backend/internal/infrastructure/logger"
)

type MQTTSubscription struct {
	mqttClient             *Client
	monitoringService      *service.MonitoringService
	installationRepository postgres.InstallationRepository
	db                     database.Database
}

func NewMQTTSubscription(
	mqttClient *Client,
	monitoringService *service.MonitoringService,
	installationRepository postgres.InstallationRepository,
	db database.Database,
) *MQTTSubscription {
	return &MQTTSubscription{
		mqttClient:             mqttClient,
		monitoringService:      monitoringService,
		installationRepository: installationRepository,
		db:                     db,
	}
}

func (subscription *MQTTSubscription) SetupMQTTSubscriptions() {
	panelIDs := subscription.installationRepository.FindAllPanelsID(subscription.db)
	for _, panelID := range panelIDs {
		subscription.mqttClient.Subscribe(fmt.Sprintf("panel/%d/status", panelID), subscription.monitoringService.HandleStatusMessage)
		subscription.mqttClient.Subscribe(fmt.Sprintf("panel/%d/history", panelID), subscription.monitoringService.HandleHistoryMessage)
		subscription.mqttClient.Subscribe(fmt.Sprintf("panel/%d/event", panelID), subscription.monitoringService.HandleEventMessage)
	}
}

func (subscription *MQTTSubscription) RefreshMQTTSubscriptions() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if subscription.mqttClient.client.IsConnected() {
			subscription.SetupMQTTSubscriptions()
		} else {
			loggerImpl.GetLogger().Warn("Skipping MQTT subscription refresh - client not connected")
		}
	}
}
