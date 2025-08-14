package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"
	mqttdto "github.com/BargheNo/Backend/internal/application/dto/mqtt"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/logger"
	"github.com/BargheNo/Backend/internal/domain/mqtt"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	loggerImpl "github.com/BargheNo/Backend/internal/infrastructure/logger"
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
			// Subscribe to specific message type topics
			service.mqttClient.Subscribe(fmt.Sprintf("panel/%d/status", panelID), service.HandleStatusMessage)
			service.mqttClient.Subscribe(fmt.Sprintf("panel/%d/history", panelID), service.HandleHistoryMessage)
			service.mqttClient.Subscribe(fmt.Sprintf("panel/%d/event", panelID), service.HandleEventMessage)
		}
	}()

	return service
}

func (monitoringService *MonitoringService) HandleStatusMessage(topic string, payload []byte) {
	panelID, err := strconv.ParseUint(topic[len("panel/"):strings.Index(topic, "/status")], 10, 32)
	if err != nil {
		loggerImpl.GetLogger().Error("Failed to parse panel ID from topic", logger.Error("error:", err))
		return
	}

	panel, err := monitoringService.installationService.GetPanelByID(uint(panelID))
	if err != nil {
		loggerImpl.GetLogger().Error("Failed to get panel", logger.Error("error:", err))
		return
	}

	// Send to websocket
	monitoringService.hub.SendToUser(panel.Customer.ID, websocket.MessageTypeMonitoring, payload)

	// Handle the status message
	monitoringService.handleStatusMessage(uint(panelID), payload)
}

func (monitoringService *MonitoringService) HandleHistoryMessage(topic string, payload []byte) {
	panelID, err := strconv.ParseUint(topic[len("panel/"):strings.Index(topic, "/history")], 10, 32)
	if err != nil {
		loggerImpl.GetLogger().Error("Failed to parse panel ID from topic", logger.Error("error:", err))
		return
	}

	panel, err := monitoringService.installationService.GetPanelByID(uint(panelID))
	if err != nil {
		loggerImpl.GetLogger().Error("Failed to get panel", logger.Error("error:", err))
		return
	}

	// Send to websocket
	monitoringService.hub.SendToUser(panel.Customer.ID, websocket.MessageTypeMonitoring, payload)

	// Handle the history message
	monitoringService.handleHistoryMessage(uint(panelID), payload)
}

func (monitoringService *MonitoringService) HandleEventMessage(topic string, payload []byte) {
	panelID, err := strconv.ParseUint(topic[len("panel/"):strings.Index(topic, "/event")], 10, 32)
	if err != nil {
		loggerImpl.GetLogger().Error("Failed to parse panel ID from topic", logger.Error("error:", err))
		return
	}

	panel, err := monitoringService.installationService.GetPanelByID(uint(panelID))
	if err != nil {
		loggerImpl.GetLogger().Error("Failed to get panel", logger.Error("error:", err))
		return
	}

	// Send to websocket
	monitoringService.hub.SendToUser(panel.Customer.ID, websocket.MessageTypeMonitoring, payload)

	// Handle the event message
	monitoringService.handleEventMessage(uint(panelID), payload)
}

func (monitoringService *MonitoringService) handleStatusMessage(panelID uint, payload []byte) {
	var mqttMsg mqttdto.StatusMessage

	if err := json.Unmarshal(payload, &mqttMsg); err != nil {
		loggerImpl.GetLogger().Error("Failed to unmarshal status message", logger.Error("error:", err))
		return
	}

	panelStatus := &entity.PanelStatus{
		DatalogSerial: mqttMsg.DatalogSerial,
		PVSerial:      mqttMsg.PVSerial,
		PVStatus:      mqttMsg.PVStatus,
		PVPowerIn:     mqttMsg.PVPowerIn,
		PV1Voltage:    mqttMsg.PV1Voltage,
		PV1Current:    mqttMsg.PV1Current,
		PV2Voltage:    mqttMsg.PV2Voltage,
		PV2Current:    mqttMsg.PV2Current,
		PVPowerOut:    mqttMsg.PVPowerOut,
		ACFreq:        mqttMsg.ACFreq,
		ACVoltage:     mqttMsg.ACVoltage,
		ACOutputPower: mqttMsg.ACOutputPower,
		Temperature:   mqttMsg.Temperature,
		BatVoltage:    mqttMsg.BatVoltage,
		BatCurrent:    mqttMsg.BatCurrent,
		BatPower:      mqttMsg.BatPower,
		GridExport:    mqttMsg.GridExport,
		GridImport:    mqttMsg.GridImport,
		PanelID:       panelID,
	}

	if err := monitoringService.monitoringRepository.CreatePanelStatus(monitoringService.db, panelStatus); err != nil {
		loggerImpl.GetLogger().Error("Failed to create panel status", logger.Error("error:", err))
		return
	}
}

func (monitoringService *MonitoringService) handleHistoryMessage(panelID uint, payload []byte) {
	var mqttMsg mqttdto.HistoryMessage

	if err := json.Unmarshal(payload, &mqttMsg); err != nil {
		loggerImpl.GetLogger().Error("Failed to unmarshal history message", logger.Error("error:", err))
		return
	}

	panelHistory := &entity.PanelHistory{
		DatalogSerial: mqttMsg.DatalogSerial,
		PVSerial:      mqttMsg.PVSerial,
		Date:          mqttMsg.Date,
		EnergyToday:   mqttMsg.EnergyToday,
		EnergyTotal:   mqttMsg.EnergyTotal,
		PanelID:       panelID,
	}

	if err := monitoringService.monitoringRepository.CreatePanelHistory(monitoringService.db, panelHistory); err != nil {
		loggerImpl.GetLogger().Error("Failed to create panel history", logger.Error("error:", err))
		return
	}
}

func (monitoringService *MonitoringService) handleEventMessage(panelID uint, payload []byte) {
	var mqttMsg mqttdto.EventMessage

	if err := json.Unmarshal(payload, &mqttMsg); err != nil {
		loggerImpl.GetLogger().Error("Failed to unmarshal event message", logger.Error("error:", err))
		return
	}

	panelEvent := &entity.PanelEvent{
		DatalogSerial: mqttMsg.DatalogSerial,
		PVSerial:      mqttMsg.PVSerial,
		EventCode:     mqttMsg.EventCode,
		Description:   mqttMsg.Description,
		Severity:      mqttMsg.Severity,
		PanelID:       panelID,
	}

	if err := monitoringService.monitoringRepository.CreatePanelEvent(monitoringService.db, panelEvent); err != nil {
		loggerImpl.GetLogger().Error("Failed to create panel event", logger.Error("error:", err))
		return
	}
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
			Timestamp:     status.CreatedAt,
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
			Timestamp:     history.CreatedAt,
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
			Timestamp:     event.CreatedAt,
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
			Timestamp:     status.CreatedAt,
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
			Timestamp:     history.CreatedAt,
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
			Timestamp:     event.CreatedAt,
		}
	}

	count, err := monitoringService.monitoringRepository.CountPanelEventByPanelID(monitoringService.db, listInfo.PanelID)
	if err != nil {
		return nil, 0, err
	}

	return response, count, nil
}
