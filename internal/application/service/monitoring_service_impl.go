package serviceimpl

import (
	"fmt"

	"github.com/BargheNo/Backend/internal/domain/mqtt"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MonitoringService struct {
	mqttClient             mqtt.Client
	installationRepository repository.InstallationRepository
	db                     database.Database
}

func NewMonitoringService(mqttClient mqtt.Client, db database.Database, installationRepository repository.InstallationRepository) *MonitoringService {
	service := &MonitoringService{mqttClient: mqttClient, db: db, installationRepository: installationRepository}
	go func() {
		panelIDs := service.installationRepository.FindAllPanelsID(service.db)
		println(len(panelIDs))
		for _, panelID := range panelIDs {
			service.mqttClient.Subscribe(fmt.Sprintf("panel/%d", panelID), service.HandleMessage)
		}
	}()

	return service
}

func (s *MonitoringService) HandleMessage(topic string, payload []byte) {
	fmt.Println(topic, string(payload))
}
