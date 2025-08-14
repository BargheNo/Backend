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

	mqttdto "github.com/BargheNo/Backend/internal/application/dto/mqtt"
	"github.com/BargheNo/Backend/internal/domain/logger"
	loggerImpl "github.com/BargheNo/Backend/internal/infrastructure/logger"
	mqttClient "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	BrokerURL string
	ClientID  string
	Topic     string
	Username  string
	Password  string
	QoS       byte
	Retained  bool
}

type Panel struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	DatalogSerial        string `json:"datalog_serial"`
	PVSerial             string `json:"pv_serial"`
	BuildingType         string `json:"building_type"`
	Area                 uint   `json:"area"`
	Power                uint   `json:"power"`
	Tilt                 uint   `json:"tilt"`
	Azimuth              uint   `json:"azimuth"`
	TotalNumberOfModules uint   `json:"total_number_of_modules"`
}

type MQTTClient struct {
	client       mqttClient.Client
	config       Config
	panel        Panel
	energyToday  float64
	energyTotal  float64
	statusTicker *time.Ticker
	historyTimer *time.Timer
	eventTimer   *time.Timer
	stopChan     chan struct{}
}

func NewMQTTClient(config Config) *MQTTClient {
	return &MQTTClient{
		config: config,
		panel: Panel{
			ID:                   1,
			Name:                 "Panel-001",
			DatalogSerial:        "YUZ081920C",
			PVSerial:             "4FZG821037",
			BuildingType:         "residential",
			Area:                 50,
			Power:                5000,
			Tilt:                 30,
			Azimuth:              180,
			TotalNumberOfModules: 16,
		},
		energyToday: 5.24,
		energyTotal: 1456.72,
		stopChan:    make(chan struct{}),
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
	}
	opts.SetTLSConfig(tlsConfig)

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
	if m.eventTimer != nil {
		m.eventTimer.Stop()
	}

	if m.client.IsConnected() {
		m.client.Disconnect(250)
		loggerImpl.GetLogger().Info("Disconnected from MQTT broker")
	}
}

func (m *MQTTClient) publishMessageWithType(payload interface{}, messageType string) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Use topic-based routing for message type identification
	topic := fmt.Sprintf("%s/%s", m.config.Topic, messageType)

	token := m.client.Publish(topic, m.config.QoS, m.config.Retained, jsonPayload)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to publish %s message: %v", messageType, token.Error())
	}

	loggerImpl.GetLogger().Info(fmt.Sprintf("%s message published", messageType), logger.String("topic", topic))
	return nil
}

func (m *MQTTClient) createStatusMessage() mqttdto.StatusMessage {
	pvPowerIn := 4200.5 + (rand.Float64()-0.5)*100
	temperature := 44.8 + (rand.Float64()-0.5)*5
	gridExport := 3900.0 + (rand.Float64()-0.5)*200

	return mqttdto.StatusMessage{
		DatalogSerial: m.panel.DatalogSerial,
		PVSerial:      m.panel.PVSerial,
		PVStatus:      1, // Active status
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

func (m *MQTTClient) createHistoryMessage() mqttdto.HistoryMessage {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	return mqttdto.HistoryMessage{
		DatalogSerial: m.panel.DatalogSerial,
		PVSerial:      m.panel.PVSerial,
		Date:          yesterday.Format("2006-01-02"),
		EnergyToday:   4.21,
		EnergyTotal:   m.energyTotal - 5.24,
	}
}

func (m *MQTTClient) createEventMessage() mqttdto.EventMessage {
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

	return mqttdto.EventMessage{
		DatalogSerial: m.panel.DatalogSerial,
		PVSerial:      m.panel.PVSerial,
		EventCode:     selectedEvent.code,
		Description:   selectedEvent.description,
		Severity:      selectedEvent.severity,
	}
}

func (m *MQTTClient) scheduleRandomEvent() {
	minDuration := 10 * time.Minute
	maxDuration := 48 * time.Hour

	randomDuration := minDuration + time.Duration(rand.Int63n(int64(maxDuration-minDuration)))

	m.eventTimer = time.AfterFunc(randomDuration, func() {
		select {
		case <-m.stopChan:
			return
		default:
			if err := m.publishMessageWithType(m.createEventMessage(), "event"); err != nil {
				loggerImpl.GetLogger().Error("Error publishing event message", logger.Error("error:", err))
			}
			m.scheduleRandomEvent()
		}
	})

	loggerImpl.GetLogger().Info("Next event scheduled", logger.Duration("duration", randomDuration))
}

func (m *MQTTClient) getNextMidnight() time.Duration {
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return nextMidnight.Sub(now)
}

func (m *MQTTClient) StartPublishers() {
	m.statusTicker = time.NewTicker(10 * time.Second)
	go func() {
		if err := m.publishMessageWithType(m.createStatusMessage(), "status"); err != nil {
			loggerImpl.GetLogger().Error("Error publishing initial status message", logger.Error("error:", err))
		}

		for {
			select {
			case <-m.statusTicker.C:
				if err := m.publishMessageWithType(m.createStatusMessage(), "status"); err != nil {
					loggerImpl.GetLogger().Error("Error publishing status message", logger.Error("error:", err))
				}
			case <-m.stopChan:
				return
			}
		}
	}()

	go func() {
		timeToMidnight := m.getNextMidnight()
		m.historyTimer = time.AfterFunc(timeToMidnight, func() {
			if err := m.publishMessageWithType(m.createHistoryMessage(), "history"); err != nil {
				loggerImpl.GetLogger().Error("Error publishing history message", logger.Error("error:", err))
			}

			dailyTicker := time.NewTicker(24 * time.Hour)
			go func() {
				defer dailyTicker.Stop()
				for {
					select {
					case <-dailyTicker.C:
						if err := m.publishMessageWithType(m.createHistoryMessage(), "history"); err != nil {
							loggerImpl.GetLogger().Error("Error publishing daily history message", logger.Error("error:", err))
						}
					case <-m.stopChan:
						return
					}
				}
			}()
		})
		loggerImpl.GetLogger().Info("First history message scheduled", logger.Duration("duration", timeToMidnight))
	}()

	m.scheduleRandomEvent()
	loggerImpl.GetLogger().Info("All publishers started")
}

var messagePubHandler mqttClient.MessageHandler = func(client mqttClient.Client, msg mqttClient.Message) {
	loggerImpl.GetLogger().Info("Received message", logger.String("payload", string(msg.Payload())), logger.String("topic", msg.Topic()))
}

func main() {
	config := Config{
		BrokerURL: "tls://e3e602004dad46329bd6b02c927ed397.s1.eu.hivemq.cloud:8883",
		ClientID:  "barghe-no-backend-sender",
		Topic:     "panel/1",
		Username:  "KianYari",
		Password:  "BargheNoMqtt4032",
		QoS:       0,
		Retained:  false,
	}

	mqttClient := NewMQTTClient(config)

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
