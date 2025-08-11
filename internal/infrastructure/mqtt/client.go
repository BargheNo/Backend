package mqttimpl

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client mqtt.Client
	config *bootstrap.MQTT
}

func NewClient(config *bootstrap.MQTT) *Client {
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%s", config.Broker, config.Port))
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)

	// Add connection stability settings
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Minute)
	opts.SetConnectRetryInterval(30 * time.Second)
	opts.SetCleanSession(true)

	// Add TLS config for HiveMQ Cloud
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Set to true only for testing
	}
	opts.SetTLSConfig(tlsConfig)

	// Connection lost handler
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v. Will auto-reconnect...", err)
	})

	// On connect handler
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("MQTT connected to broker")
	})

	// Reconnect handler
	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {
		log.Println("MQTT attempting to reconnect...")
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &Client{client: client, config: config}
}

func (c *Client) Subscribe(topic string, handler func(topic string, payload []byte)) {
	token := c.client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Topic(), msg.Payload())
	})
	if token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe to topic %s: %v", topic, token.Error())
		return
	}
	log.Printf("Successfully subscribed to topic: %s", topic)
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
}
