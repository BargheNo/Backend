package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	mqttdto "github.com/BargheNo/Backend/internal/application/dto/mqtt"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/logger"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	loggerImpl "github.com/BargheNo/Backend/internal/infrastructure/logger"
	infraPostgres "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
	mqttClient "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	BrokerURL        string
	ClientID         string
	Topic            string
	Username         string
	Password         string
	QoS              byte
	Retained         bool
	StatusInterval   time.Duration
	HistoryInterval  time.Duration
	EventMinInterval time.Duration
	EventMaxInterval time.Duration
}

type MQTTClient struct {
	client           mqttClient.Client
	config           Config
	db               database.Database
	installationRepo postgres.InstallationRepository
	panels           []*entity.Panel
	energyToday      map[uint]float64
	energyTotal      map[uint]float64
	statusTicker     *time.Ticker
	historyTimer     *time.Timer
	eventTimers      map[uint]*time.Timer
	stopChan         chan struct{}
}

func NewMQTTClient(config Config, db database.Database) *MQTTClient {
	return &MQTTClient{
		config:           config,
		db:               db,
		installationRepo: infraPostgres.NewInstallationRepository(),
		energyToday:      make(map[uint]float64),
		energyTotal:      make(map[uint]float64),
		stopChan:         make(chan struct{}),
		eventTimers:      make(map[uint]*time.Timer),
	}
}

func (m *MQTTClient) Connect() error {
	opts := mqttClient.NewClientOptions()
	opts.AddBroker(m.config.BrokerURL)
	opts.SetClientID(m.config.ClientID)
	opts.SetUsername(m.config.Username)
	opts.SetPassword(m.config.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}
	opts.SetTLSConfig(tlsConfig)

	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Minute)
	opts.SetCleanSession(true)

	opts.SetConnectionLostHandler(func(client mqttClient.Client, err error) {
		loggerImpl.GetLogger().Error("MQTT connection lost", logger.Error("error:", err))
	})

	opts.SetOnConnectHandler(func(client mqttClient.Client) {
		loggerImpl.GetLogger().Info("Connected to MQTT broker")
	})

	m.client = mqttClient.NewClient(opts)

	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	return nil
}

func (m *MQTTClient) Disconnect() {
	close(m.stopChan)

	if m.statusTicker != nil {
		m.statusTicker.Stop()
	}
	if m.historyTimer != nil {
		m.historyTimer.Stop()
	}

	for _, timer := range m.eventTimers {
		if timer != nil {
			timer.Stop()
		}
	}

	if m.client.IsConnected() {
		m.client.Disconnect(250)
		loggerImpl.GetLogger().Info("Disconnected from MQTT broker")
	}
}

func (m *MQTTClient) LoadPanelsFromDatabase() error {
	allowedStatuses := []enum.PanelStatus{
		enum.PanelStatusActive,
	}

	panels, err := m.installationRepo.FindPanelsByStatus(m.db, allowedStatuses, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch panels from database: %v", err)
	}

	if len(panels) == 0 {
		loggerImpl.GetLogger().Warn("No active panels found in database")
		return nil
	}

	m.panels = panels

	for _, panel := range panels {
		m.energyToday[panel.ID] = 4.0 + rand.Float64()*3.0
		m.energyTotal[panel.ID] = 1000.0 + rand.Float64()*1000.0
	}

	loggerImpl.GetLogger().Info(fmt.Sprintf("Loaded %d panels from database", len(panels)))
	return nil
}

func (m *MQTTClient) publishMessageWithType(payload interface{}, messageType string, panelID uint) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	topic := fmt.Sprintf("%s/%d/%s", m.config.Topic, panelID, messageType)

	token := m.client.Publish(topic, m.config.QoS, m.config.Retained, jsonPayload)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to publish %s message for panel %d: %v", messageType, panelID, token.Error())
	}

	loggerImpl.GetLogger().Info(fmt.Sprintf("%s message published for panel %d", messageType, panelID), logger.String("topic", topic))
	return nil
}

func (m *MQTTClient) createStatusMessage(panel *entity.Panel) mqttdto.StatusMessage {
	pvPowerIn := 4200.5 + (rand.Float64()-0.5)*100
	temperature := 44.8 + (rand.Float64()-0.5)*5
	gridExport := 3900.0 + (rand.Float64()-0.5)*200

	datalogSerial := fmt.Sprintf("YUZ%08d", panel.ID)
	pvSerial := fmt.Sprintf("4FZG%08d", panel.ID)

	return mqttdto.StatusMessage{
		DatalogSerial: datalogSerial,
		PVSerial:      pvSerial,
		PVStatus:      1,
		PVPowerIn:     pvPowerIn,
		PV1Voltage:    350.0,
		PV1Current:    6.0,
		PV2Voltage:    348.5,
		PV2Current:    6.2,
		PVPowerOut:    4150.0,
		ACFreq:        50.0,
		ACVoltage:     230.4,
		ACOutputPower: 4120.0,
		Temperature:   temperature,
		BatVoltage:    51.2,
		BatCurrent:    -5.3,
		BatPower:      -270.0,
		GridExport:    gridExport,
		GridImport:    0.0,
	}
}

func (m *MQTTClient) createHistoryMessage(panel *entity.Panel) mqttdto.HistoryMessage {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	datalogSerial := fmt.Sprintf("YUZ%08d", panel.ID)
	pvSerial := fmt.Sprintf("4FZG%08d", panel.ID)

	return mqttdto.HistoryMessage{
		DatalogSerial: datalogSerial,
		PVSerial:      pvSerial,
		Date:          yesterday.Format("2006-01-02"),
		EnergyToday:   m.energyToday[panel.ID],
		EnergyTotal:   m.energyTotal[panel.ID],
	}
}

func (m *MQTTClient) createEventMessage(panel *entity.Panel) mqttdto.EventMessage {
	events := []struct {
		code        string
		description string
		severity    string
	}{
		{code: "P28", description: "PV isolation fault", severity: "error"},
		{code: "P01", description: "Grid overvoltage", severity: "warning"},
		{code: "P02", description: "Grid undervoltage", severity: "warning"},
		{code: "P05", description: "Temperature high", severity: "warning"},
		{code: "P15", description: "Battery low voltage", severity: "info"},
	}

	selectedEvent := events[rand.Intn(len(events))]

	datalogSerial := fmt.Sprintf("YUZ%08d", panel.ID)
	pvSerial := fmt.Sprintf("4FZG%08d", panel.ID)

	return mqttdto.EventMessage{
		DatalogSerial: datalogSerial,
		PVSerial:      pvSerial,
		EventCode:     selectedEvent.code,
		Description:   selectedEvent.description,
		Severity:      selectedEvent.severity,
	}
}

func (m *MQTTClient) scheduleRandomEvent(panel *entity.Panel) {
	minDuration := m.config.EventMinInterval
	maxDuration := m.config.EventMaxInterval

	randomDuration := minDuration + time.Duration(rand.Int63n(int64(maxDuration-minDuration)))

	m.eventTimers[panel.ID] = time.AfterFunc(randomDuration, func() {
		select {
		case <-m.stopChan:
			return
		default:
			if err := m.publishMessageWithType(m.createEventMessage(panel), "event", panel.ID); err != nil {
				loggerImpl.GetLogger().Error("Error publishing event message", logger.Error("error:", err))
			}
			m.scheduleRandomEvent(panel)
		}
	})

	loggerImpl.GetLogger().Info(fmt.Sprintf("Next event scheduled for panel %d", panel.ID), logger.Duration("duration", randomDuration))
}

func (m *MQTTClient) getNextMidnight() time.Duration {
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return nextMidnight.Sub(now)
}

func (m *MQTTClient) StartPublishers() {
	if err := m.LoadPanelsFromDatabase(); err != nil {
		loggerImpl.GetLogger().Error("Failed to load panels from database", logger.Error("error:", err))
		return
	}

	if len(m.panels) == 0 {
		loggerImpl.GetLogger().Warn("No panels to publish data for")
		return
	}

	m.statusTicker = time.NewTicker(m.config.StatusInterval)
	go func() {
		for _, panel := range m.panels {
			if err := m.publishMessageWithType(m.createStatusMessage(panel), "status", panel.ID); err != nil {
				loggerImpl.GetLogger().Error(fmt.Sprintf("Error publishing initial status message for panel %d", panel.ID), logger.Error("error:", err))
			}
		}

		for {
			select {
			case <-m.statusTicker.C:
				for _, panel := range m.panels {
					if err := m.publishMessageWithType(m.createStatusMessage(panel), "status", panel.ID); err != nil {
						loggerImpl.GetLogger().Error(fmt.Sprintf("Error publishing status message for panel %d", panel.ID), logger.Error("error:", err))
					}
				}
			case <-m.stopChan:
				return
			}
		}
	}()

	go func() {
		timeToMidnight := m.getNextMidnight()
		m.historyTimer = time.AfterFunc(timeToMidnight, func() {
			for _, panel := range m.panels {
				if err := m.publishMessageWithType(m.createHistoryMessage(panel), "history", panel.ID); err != nil {
					loggerImpl.GetLogger().Error(fmt.Sprintf("Error publishing history message for panel %d", panel.ID), logger.Error("error:", err))
				}
			}

			dailyTicker := time.NewTicker(m.config.HistoryInterval)
			go func() {
				defer dailyTicker.Stop()
				for {
					select {
					case <-dailyTicker.C:
						for _, panel := range m.panels {
							if err := m.publishMessageWithType(m.createHistoryMessage(panel), "history", panel.ID); err != nil {
								loggerImpl.GetLogger().Error(fmt.Sprintf("Error publishing daily history message for panel %d", panel.ID), logger.Error("error:", err))
							}
						}
					case <-m.stopChan:
						return
					}
				}
			}()
		})
		loggerImpl.GetLogger().Info("First history message scheduled", logger.Duration("duration", timeToMidnight))
	}()

	for _, panel := range m.panels {
		m.scheduleRandomEvent(panel)
	}

	panelReloadTicker := time.NewTicker(20 * time.Second)
	go func() {
		defer panelReloadTicker.Stop()
		for {
			select {
			case <-panelReloadTicker.C:
				if err := m.LoadPanelsFromDatabase(); err != nil {
					loggerImpl.GetLogger().Error("Failed to reload panels from database", logger.Error("error:", err))
				} else {
					loggerImpl.GetLogger().Info("Panels reloaded from database")
				}
			case <-m.stopChan:
				return
			}
		}
	}()

	loggerImpl.GetLogger().Info(fmt.Sprintf("All publishers started for %d panels", len(m.panels)))
}

var messagePubHandler mqttClient.MessageHandler = func(client mqttClient.Client, msg mqttClient.Message) {
	loggerImpl.GetLogger().Info("Received message", logger.String("payload", string(msg.Payload())), logger.String("topic", msg.Topic()))
}

func main() {
	config := bootstrap.Run()

	mqttConfig := Config{
		BrokerURL:        fmt.Sprintf("tls://%s:%s", config.Env.MQTT.Broker, config.Env.MQTT.Port),
		ClientID:         config.Env.MQTT.SenderClientID,
		Topic:            "panel",
		Username:         config.Env.MQTT.Username,
		Password:         config.Env.MQTT.Password,
		QoS:              0,
		Retained:         false,
		StatusInterval:   30 * time.Second,
		HistoryInterval:  24 * time.Hour,
		EventMinInterval: 30 * time.Minute,
		EventMaxInterval: 6 * time.Hour,
	}

	db := database.NewPostgresDatabase(&config.Env.PrimaryDB)

	mqttClient := NewMQTTClient(mqttConfig, db)

	if err := mqttClient.Connect(); err != nil {
		loggerImpl.GetLogger().Error("Failed to connect to MQTT broker", logger.Error("error:", err))
		os.Exit(1)
	}
	defer mqttClient.Disconnect()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	mqttClient.StartPublishers()

	sig := <-sigChan
	loggerImpl.GetLogger().Info("Received signal, shutting down gracefully", logger.String("signal", sig.String()))
}
