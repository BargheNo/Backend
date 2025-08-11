package service

import (
	"fmt"
	"strconv"

	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/mqtt"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
)

type MonitoringService struct {
	mqttClient             mqtt.Client
	installationService    usecase.InstallationService
	installationRepository repository.InstallationRepository
	db                     database.Database
	hub                    *websocket.Hub
}

func NewMonitoringService(
	mqttClient mqtt.Client,
	db database.Database,
	installationRepository repository.InstallationRepository,
	hub *websocket.Hub,
	installationService usecase.InstallationService,
) *MonitoringService {
	service := &MonitoringService{
		mqttClient:             mqttClient,
		db:                     db,
		installationRepository: installationRepository,
		hub:                    hub,
		installationService:    installationService,
	}
	go func() {
		panelIDs := service.installationRepository.FindAllPanelsID(service.db)
		for _, panelID := range panelIDs {
			service.mqttClient.Subscribe(fmt.Sprintf("panel/%d", panelID), service.HandleMessage)
		}
	}()

	return service
}

func (s *MonitoringService) HandleMessage(topic string, payload []byte) {
	panelID, err := strconv.ParseUint(topic[len("panel/"):], 10, 32)
	if err != nil {
		panic(err)
	}
	panel, err := s.installationService.GetPanelByID(uint(panelID))
	if err != nil {
		panic(err)
	}

	s.hub.SendToUser(panel.Customer.ID, websocket.MessageTypeMonitoring, payload)
}
