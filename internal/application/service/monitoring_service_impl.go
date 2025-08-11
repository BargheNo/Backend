package service

import (
	"fmt"
	"strconv"

	monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"
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
	monitoringRepository   repository.MonitoringRepository
	db                     database.Database
	hub                    *websocket.Hub
}

func NewMonitoringService(
	mqttClient mqtt.Client,
	db database.Database,
	installationRepository repository.InstallationRepository,
	monitoringRepository repository.MonitoringRepository,
	hub *websocket.Hub,
	installationService usecase.InstallationService,
) *MonitoringService {
	service := &MonitoringService{
		mqttClient:             mqttClient,
		db:                     db,
		installationRepository: installationRepository,
		monitoringRepository:   monitoringRepository,
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

func (s *MonitoringService) GetPanelStatus(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.CustomerPanelStatusResponse, int64, error) {
	_, err := s.installationService.ValidatePanelOwnership(listInfo.PanelID, listInfo.OwnerID)
	if err != nil {
		return nil, 0, err
	}

	panelStatus, err := s.monitoringRepository.FindPanelStatusByPanelID(s.db, listInfo.PanelID)
	if err != nil {
		return nil, 0, err
	}

	response := make([]monitoringdto.CustomerPanelStatusResponse, len(panelStatus))
	for i, status := range panelStatus {
		response[i] = monitoringdto.CustomerPanelStatusResponse{
			DatalogSerial: status.DatalogSerial,
			PVSerial:      status.PVSerial,
			PVStatus:      status.PVStatus,
			PVPowerIn:     status.PVPowerIn,
			PV1Voltage:    status.PV1Voltage,
			PV1Current:    status.PV1Current,
			PV2Voltage:    status.PV2Voltage,
			PV2Current:    status.PV2Current,
			PVPowerOut:    status.PVPowerOut,
			ACFreq:        status.ACFreq,
			ACVoltage:     status.ACVoltage,
			ACOutputPower: status.ACOutputPower,
			Temperature:   status.Temperature,
			BatVoltage:    status.BatVoltage,
			BatCurrent:    status.BatCurrent,
			BatPower:      status.BatPower,
			GridExport:    status.GridExport,
			GridImport:    status.GridImport,
			EnergyToday:   status.EnergyToday,
			EnergyTotal:   status.EnergyTotal,
			Timestamp:     status.Timestamp,
		}
	}
	return response, int64(len(response)), nil
}
