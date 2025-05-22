package mqttimpl

import (
	"fmt"

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

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &Client{client: client, config: config}
}

func (c *Client) Subscribe(topic string) {
	token := c.client.Subscribe(topic, 1, nil)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
}
