package service

import (
	"fmt"
	"strconv"

	monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/mqtt"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
)

type MonitoringService struct {
	mqttClient             mqtt.Client
	installationService    usecase.InstallationService
	installationRepository postgres.InstallationRepository
	monitoringRepository   postgres.MonitoringRepository
	db                     database.Database
	hub                    *websocket.Hub
}

func NewMonitoringService(
	mqttClient mqtt.Client,
	db database.Database,
	installationRepository postgres.InstallationRepository,
	monitoringRepository postgres.MonitoringRepository,
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

func (s *MonitoringService) GetPanelStatus(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelStatusResponse, int64, error) {
	_, err := s.installationService.ValidatePanelOwnership(listInfo.PanelID, listInfo.OwnerID)
	if err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().WithPagination(listInfo.Limit, listInfo.Offset)

	panelStatus, err := s.monitoringRepository.FindPanelStatusByPanelID(s.db, listInfo.PanelID, options)
	if err != nil {
		return nil, 0, err
	}

	response := make([]monitoringdto.PanelStatusResponse, len(panelStatus))
	for i, status := range panelStatus {
		response[i] = monitoringdto.PanelStatusResponse{
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

	count, err := s.monitoringRepository.CountPanelStatusByPanelID(s.db, listInfo.PanelID)
	if err != nil {
		return nil, 0, err
	}

	return response, count, nil
}

func (monitoringService *MonitoringService) GetPanelHistory(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelHistoryResponse, int64, error) {
	_, err := monitoringService.installationService.ValidatePanelOwnership(listInfo.PanelID, listInfo.OwnerID)
	if err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().WithPagination(listInfo.Limit, listInfo.Offset)

	panelHistory, err := monitoringService.monitoringRepository.FindPanelHistoryByPanelID(monitoringService.db, listInfo.PanelID, options)
	if err != nil {
		return nil, 0, err
	}

	response := make([]monitoringdto.PanelHistoryResponse, len(panelHistory))
	for i, history := range panelHistory {
		response[i] = monitoringdto.PanelHistoryResponse{
			DatalogSerial: history.DatalogSerial,
			PVSerial:      history.PVSerial,
			Date:          history.Date,
			EnergyToday:   history.EnergyToday,
			EnergyTotal:   history.EnergyTotal,
			Timestamp:     history.Timestamp,
		}
	}

	count, err := monitoringService.monitoringRepository.CountPanelHistoryByPanelID(monitoringService.db, listInfo.PanelID)
	if err != nil {
		return nil, 0, err
	}
	return response, count, nil
}

func (monitoringService *MonitoringService) GetPanelEvent(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelEventResponse, int64, error) {
	_, err := monitoringService.installationService.ValidatePanelOwnership(listInfo.PanelID, listInfo.OwnerID)
	if err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().WithPagination(listInfo.Limit, listInfo.Offset)

	panelEvent, err := monitoringService.monitoringRepository.FindPanelEventByPanelID(monitoringService.db, listInfo.PanelID, options)
	if err != nil {
		return nil, 0, err
	}

	response := make([]monitoringdto.PanelEventResponse, len(panelEvent))
	for i, event := range panelEvent {
		response[i] = monitoringdto.PanelEventResponse{
			DatalogSerial: event.DatalogSerial,
			PVSerial:      event.PVSerial,
			EventCode:     event.EventCode,
			Description:   event.Description,
			Severity:      event.Severity,
			Timestamp:     event.Timestamp,
		}
	}

	count, err := monitoringService.monitoringRepository.CountPanelEventByPanelID(monitoringService.db, listInfo.PanelID)
	if err != nil {
		return nil, 0, err
	}

	return response, count, nil
}
