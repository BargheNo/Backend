package mqtt

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/logger"
	loggerImpl "github.com/BargheNo/Backend/internal/infrastructure/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client mqtt.Client
	config *bootstrap.MQTT
}

func NewClient(config *bootstrap.MQTT) *Client {
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		loggerImpl.GetLogger().Info("Received message", logger.String("payload", string(msg.Payload())), logger.String("topic", msg.Topic()))
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%s", config.Broker, config.Port))
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)

	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Minute)
	opts.SetConnectRetryInterval(30 * time.Second)
	opts.SetCleanSession(false)
	opts.SetResumeSubs(true)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}
	opts.SetTLSConfig(tlsConfig)

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		loggerImpl.GetLogger().Warn("MQTT connection lost", logger.Error("error", err))
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		loggerImpl.GetLogger().Info("MQTT reconnected successfully")
	})

	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &Client{client: client, config: config}
}

func (c *Client) Subscribe(topic string, handler func(topic string, payload []byte)) {
	if !c.client.IsConnected() {
		loggerImpl.GetLogger().Warn("Cannot subscribe: MQTT client not connected", logger.String("topic", topic))
		return
	}

	token := c.client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Topic(), msg.Payload())
	})
	if token.Wait() && token.Error() != nil {
		loggerImpl.GetLogger().Error("Failed to subscribe to topic", logger.Error("error:", token.Error()), logger.String("topic", topic))
		return
	}

	loggerImpl.GetLogger().Info("Successfully subscribed to topic", logger.String("topic", topic))
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
}
