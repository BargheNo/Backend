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
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("ssl://%s:%s", config.Broker, config.Port)).
		SetClientID(config.ClientID).
		SetUsername(config.Username).
		SetPassword(config.Password)

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &Client{client: c, config: config}
}

func (c *Client) Subscribe(topic string, handler func(payload []byte)) {
	token := c.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Payload())
		print(msg.Payload(), "\n")
	})
	print(2, "\n")
	print(token.Error(), "\n")
	if token.Wait() && token.Error() != nil {
		// panic(token.Error())
		// fmt.Println(token.Error())
		print(1, "\n", token.Error())
	}
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
}
