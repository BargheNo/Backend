package service

import (
	"fmt"
	"strconv"

	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
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
	corporationService     usecase.CorporationService
	installationRepository postgres.InstallationRepository
	monitoringRepository   postgres.MonitoringRepository
	db                     database.Database
	hub                    *websocket.Hub
}

func NewMonitoringService(
	mqttClient mqtt.Client,
	db database.Database,
	corporationService usecase.CorporationService,
	installationRepository postgres.InstallationRepository,
	monitoringRepository postgres.MonitoringRepository,
	hub *websocket.Hub,
	installationService usecase.InstallationService,
) *MonitoringService {
	service := &MonitoringService{
		mqttClient:             mqttClient,
		db:                     db,
		corporationService:     corporationService,
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

func (s *MonitoringService) GetCustomerPanelStatus(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelStatusResponse, int64, error) {
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

func (monitoringService *MonitoringService) GetCustomerPanelHistory(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelHistoryResponse, int64, error) {
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

func (monitoringService *MonitoringService) GetCustomerPanelEvent(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelEventResponse, int64, error) {
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

func (monitoringService *MonitoringService) GetCorporationPanelStatus(listInfo monitoringdto.CorporationPanelStatusListRequest) ([]monitoringdto.PanelStatusResponse, int64, error) {
	if err := monitoringService.corporationService.CheckApplicantAccess(listInfo.CorporationID, listInfo.UserID); err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().WithPagination(listInfo.Limit, listInfo.Offset)

	getPanelRequest := installationdto.CorporationPanelRequest{
		CorporationID:  listInfo.CorporationID,
		OperatorID:     listInfo.UserID,
		InstallationID: listInfo.PanelID,
	}
	_, err := monitoringService.installationService.GetCorporationPanel(getPanelRequest)
	if err != nil {
		return nil, 0, err
	}

	panelStatus, err := monitoringService.monitoringRepository.FindPanelStatusByPanelID(monitoringService.db, listInfo.PanelID, options)
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

	count, err := monitoringService.monitoringRepository.CountPanelStatusByPanelID(monitoringService.db, listInfo.PanelID)
	if err != nil {
		return nil, 0, err
	}

	return response, count, nil
}

func (monitoringService *MonitoringService) GetCorporationPanelHistory(listInfo monitoringdto.CorporationPanelStatusListRequest) ([]monitoringdto.PanelHistoryResponse, int64, error) {
	if err := monitoringService.corporationService.CheckApplicantAccess(listInfo.CorporationID, listInfo.UserID); err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().WithPagination(listInfo.Limit, listInfo.Offset)

	getPanelRequest := installationdto.CorporationPanelRequest{
		CorporationID:  listInfo.CorporationID,
		OperatorID:     listInfo.UserID,
		InstallationID: listInfo.PanelID,
	}
	_, err := monitoringService.installationService.GetCorporationPanel(getPanelRequest)
	if err != nil {
		return nil, 0, err
	}

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

func (monitoringService *MonitoringService) GetCorporationPanelEvent(listInfo monitoringdto.CorporationPanelStatusListRequest) ([]monitoringdto.PanelEventResponse, int64, error) {
	if err := monitoringService.corporationService.CheckApplicantAccess(listInfo.CorporationID, listInfo.UserID); err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().WithPagination(listInfo.Limit, listInfo.Offset)

	getPanelRequest := installationdto.CorporationPanelRequest{
		CorporationID:  listInfo.CorporationID,
		OperatorID:     listInfo.UserID,
		InstallationID: listInfo.PanelID,
	}
	_, err := monitoringService.installationService.GetCorporationPanel(getPanelRequest)
	if err != nil {
		return nil, 0, err
	}

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
