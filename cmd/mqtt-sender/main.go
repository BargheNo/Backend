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

	"github.com/BargheNo/Backend/internal/domain/logger"
	loggerImpl "github.com/BargheNo/Backend/internal/infrastructure/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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

type Device struct {
	DatalogSerial string `json:"datalogserial"`
	PVSerial      string `json:"pvserial"`
	Model         string `json:"model,omitempty"`
	Firmware      string `json:"firmware,omitempty"`
	Hardware      string `json:"hardware,omitempty"`
	LoggerFW      string `json:"logger_fw,omitempty"`
	Timezone      string `json:"timezone,omitempty"`
}

type StatusMessage struct {
	Type      string `json:"type"`
	Device    Device `json:"device"`
	Status    Status `json:"status"`
	Energy    Energy `json:"energy"`
	Timestamp string `json:"timestamp"`
}

type Status struct {
	PVStatus      int     `json:"pvstatus"`
	PVPowerIn     float64 `json:"pvpowerin"`
	PV1Voltage    float64 `json:"pv1voltage"`
	PV1Current    float64 `json:"pv1current"`
	PV2Voltage    float64 `json:"pv2voltage"`
	PV2Current    float64 `json:"pv2current"`
	PVPowerOut    float64 `json:"pvpowerout"`
	ACFreq        float64 `json:"acfreq"`
	ACVoltage     float64 `json:"acvoltage"`
	ACOutputPower float64 `json:"acoutputpower"`
	Temperature   float64 `json:"temperature"`
	BatVoltage    float64 `json:"batvoltage"`
	BatCurrent    float64 `json:"batcurrent"`
	BatPower      float64 `json:"batpower"`
	GridExport    float64 `json:"gridexport"`
	GridImport    float64 `json:"gridimport"`
}

type Energy struct {
	EnergyToday float64 `json:"energyToday"`
	EnergyTotal float64 `json:"energyTotal"`
}

type DeviceInfoMessage struct {
	Type      string `json:"type"`
	Device    Device `json:"device"`
	Timestamp string `json:"timestamp"`
}

type HistoryMessage struct {
	Type      string        `json:"type"`
	Device    Device        `json:"device"`
	History   []HistoryItem `json:"history"`
	Timestamp string        `json:"timestamp"`
}

type HistoryItem struct {
	Date        string  `json:"date"`
	EnergyToday float64 `json:"energyToday"`
	EnergyTotal float64 `json:"energyTotal"`
}

type EventMessage struct {
	Type      string `json:"type"`
	Device    Device `json:"device"`
	Event     Event  `json:"event"`
	Timestamp string `json:"timestamp"`
}

type Event struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type MQTTClient struct {
	client       mqtt.Client
	config       Config
	device       Device
	energyToday  float64
	energyTotal  float64
	statusTicker *time.Ticker
	deviceTicker *time.Ticker
	historyTimer *time.Timer
	eventTimer   *time.Timer
	stopChan     chan struct{}
}

func NewMQTTClient(config Config) *MQTTClient {
	return &MQTTClient{
		config: config,
		device: Device{
			DatalogSerial: "YUZ081920C",
			PVSerial:      "4FZG821037",
			Model:         "SPH5000TL-BL",
			Firmware:      "01.20.07",
			Hardware:      "H4.2",
			LoggerFW:      "V2.05.15",
			Timezone:      "UTC+02:00",
		},
		energyToday: 5.24,
		energyTotal: 1456.72,
		stopChan:    make(chan struct{}),
	}
}

func (m *MQTTClient) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(m.config.BrokerURL)
	opts.SetClientID(m.config.ClientID)
	opts.SetUsername(m.config.Username)
	opts.SetPassword(m.config.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	opts.SetTLSConfig(tlsConfig)

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		loggerImpl.GetLogger().Error("MQTT connection lost", logger.Error("error:", err))
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		loggerImpl.GetLogger().Info("Connected to MQTT broker")
	})

	m.client = mqtt.NewClient(opts)

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
	if m.deviceTicker != nil {
		m.deviceTicker.Stop()
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

func (m *MQTTClient) publishMessage(payload interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	token := m.client.Publish(m.config.Topic, m.config.QoS, m.config.Retained, jsonPayload)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to publish message: %v", token.Error())
	}

	loggerImpl.GetLogger().Info("Message published", logger.String("topic", m.config.Topic), logger.String("payload", string(jsonPayload)))
	return nil
}

func (m *MQTTClient) createStatusMessage() StatusMessage {
	pvPowerIn := 4200.5 + (rand.Float64()-0.5)*100
	temperature := 44.8 + (rand.Float64()-0.5)*5
	gridExport := 3900.0 + (rand.Float64()-0.5)*200

	return StatusMessage{
		Type: "status",
		Device: Device{
			DatalogSerial: m.device.DatalogSerial,
			PVSerial:      m.device.PVSerial,
		},
		Status: Status{
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
		},
		Energy: Energy{
			EnergyToday: m.energyToday,
			EnergyTotal: m.energyTotal,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func (m *MQTTClient) createDeviceInfoMessage() DeviceInfoMessage {
	return DeviceInfoMessage{
		Type:      "device_info",
		Device:    m.device,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func (m *MQTTClient) createHistoryMessage() HistoryMessage {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	dayBeforeYesterday := now.AddDate(0, 0, -2)

	return HistoryMessage{
		Type: "history",
		Device: Device{
			DatalogSerial: m.device.DatalogSerial,
			PVSerial:      m.device.PVSerial,
		},
		History: []HistoryItem{
			{
				Date:        yesterday.Format("2006-01-02"),
				EnergyToday: 4.21,
				EnergyTotal: m.energyTotal - 5.24,
			},
			{
				Date:        dayBeforeYesterday.Format("2006-01-02"),
				EnergyToday: 3.88,
				EnergyTotal: m.energyTotal - 9.45,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func (m *MQTTClient) createEventMessage() EventMessage {
	events := []Event{
		{Code: "P28", Description: "PV isolation fault", Severity: "error"},
		{Code: "P01", Description: "Grid overvoltage", Severity: "warning"},
		{Code: "P02", Description: "Grid undervoltage", Severity: "warning"},
		{Code: "P05", Description: "Temperature high", Severity: "warning"},
		{Code: "P15", Description: "Battery low voltage", Severity: "info"},
	}

	selectedEvent := events[rand.Intn(len(events))]

	return EventMessage{
		Type: "event",
		Device: Device{
			DatalogSerial: m.device.DatalogSerial,
			PVSerial:      m.device.PVSerial,
		},
		Event:     selectedEvent,
		Timestamp: time.Now().Format(time.RFC3339),
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
			if err := m.publishMessage(m.createEventMessage()); err != nil {
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
		if err := m.publishMessage(m.createStatusMessage()); err != nil {
			loggerImpl.GetLogger().Error("Error publishing initial status message", logger.Error("error:", err))
		}

		for {
			select {
			case <-m.statusTicker.C:
				if err := m.publishMessage(m.createStatusMessage()); err != nil {
					loggerImpl.GetLogger().Error("Error publishing status message", logger.Error("error:", err))
				}
			case <-m.stopChan:
				return
			}
		}
	}()

	m.deviceTicker = time.NewTicker(10 * time.Minute)
	go func() {
		time.Sleep(30 * time.Second)
		if err := m.publishMessage(m.createDeviceInfoMessage()); err != nil {
			loggerImpl.GetLogger().Error("Error publishing initial device info message", logger.Error("error:", err))
		}

		for {
			select {
			case <-m.deviceTicker.C:
				if err := m.publishMessage(m.createDeviceInfoMessage()); err != nil {
					loggerImpl.GetLogger().Error("Error publishing device info message", logger.Error("error:", err))
				}
			case <-m.stopChan:
				return
			}
		}
	}()

	go func() {
		timeToMidnight := m.getNextMidnight()
		m.historyTimer = time.AfterFunc(timeToMidnight, func() {
			if err := m.publishMessage(m.createHistoryMessage()); err != nil {
				loggerImpl.GetLogger().Error("Error publishing history message", logger.Error("error:", err))
			}

			dailyTicker := time.NewTicker(24 * time.Hour)
			go func() {
				defer dailyTicker.Stop()
				for {
					select {
					case <-dailyTicker.C:
						if err := m.publishMessage(m.createHistoryMessage()); err != nil {
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

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
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
