package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	config       *bootstrap.RabbitMQ
	constants    *bootstrap.RabbitMQConstants
	exchanges    map[string]string
	queues       map[string]bool
	bindings     map[string][]string
	isConnected  bool
	closeChannel chan struct{}
	mu           sync.RWMutex
}

func NewRabbitMQ(config *bootstrap.RabbitMQ, constants *bootstrap.RabbitMQConstants) *RabbitMQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, config.VHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		panic(err)
	}

	rmq := &RabbitMQ{
		conn:         conn,
		channel:      ch,
		config:       config,
		constants:    constants,
		exchanges:    make(map[string]string),
		queues:       make(map[string]bool),
		bindings:     make(map[string][]string),
		isConnected:  true,
		closeChannel: make(chan struct{}),
	}

	if err := rmq.declareExchange(constants.Exchange.Notifications, constants.Exchange.TypeTopic); err != nil {
		rmq.conn.Close()
		log.Printf("error during declare exchange: %v", err)
		panic(err)
	}

	if err := rmq.setupDeadLetterQueue(); err != nil {
		rmq.Close()
		log.Printf("error during declare DLQ: %v", err)
		panic(err)
	}

	queues := []string{constants.Events.UserRegistered, constants.Events.NotificationsEmail, constants.Events.NotificationsPush, constants.Events.SendNotification}

	for _, queue := range queues {
		if err := rmq.declareQueueWithDLX(queue, constants.Exchange.DLX); err != nil {
			rmq.Close()
			log.Printf("error during declare Queue: %v", err)
			panic(err)
		}
		if err := rmq.bindQueue(queue, constants.Exchange.Notifications, queue); err != nil {
			rmq.Close()
			log.Printf("error during bind Queue: %v", err)
			panic(err)
		}
	}

	go rmq.monitorConnection()

	return rmq
}

func (rmq *RabbitMQ) declareExchange(name, exchangeType string) error {
	err := rmq.channel.ExchangeDeclare(
		name,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err == nil {
		rmq.exchanges[name] = exchangeType
	}
	return err
}

func (rmq *RabbitMQ) declareQueueWithDLX(name, dlx string) error {
	args := amqp.Table{
		rmq.constants.Headers.DeadLetter: dlx,
	}

	_, err := rmq.channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		args,
	)
	if err == nil {
		rmq.queues[name] = true
	}
	return err
}

func (rmq *RabbitMQ) bindQueue(queue, exchange, routingKey string) error {
	err := rmq.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err == nil {
		rmq.bindings[queue] = append(rmq.bindings[queue], routingKey)
	}
	return err
}

func (rmq *RabbitMQ) setupDeadLetterQueue() error {
	if err := rmq.declareExchange(rmq.constants.Exchange.DLX, rmq.constants.Exchange.TypeFanout); err != nil {
		return err
	}

	_, err := rmq.channel.QueueDeclare(
		rmq.constants.Queue.DLQ,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	rmq.queues[rmq.constants.Queue.DLQ] = true

	return rmq.channel.QueueBind(
		rmq.constants.Queue.DLQ,
		"",
		rmq.constants.Exchange.DLX,
		false,
		nil,
	)
}

func (rmq *RabbitMQ) monitorConnection() {
	connCloseChan := rmq.conn.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-rmq.closeChannel:
			return
		case err := <-connCloseChan:
			if err != nil {
				rmq.mu.Lock()
				rmq.isConnected = false
				rmq.mu.Unlock()

				log.Printf("RabbitMQ connection lost: %v, attempting to reconnect...", err)

				for {
					rmq.mu.RLock()
					connected := rmq.isConnected
					rmq.mu.RUnlock()

					if connected {
						break
					}

					if err := rmq.reconnect(); err != nil {
						log.Printf("Failed to reconnect to RabbitMQ: %v, retrying in %s", err, rmq.config.RetryDelay)
						time.Sleep(rmq.config.RetryDelay)
					} else {
						log.Println("Successfully reconnected to RabbitMQ")
						connCloseChan = rmq.conn.NotifyClose(make(chan *amqp.Error))
						break
					}
				}
			}
		}
	}
}

func (rmq *RabbitMQ) reconnect() error {
	if rmq.conn != nil {
		rmq.conn.Close()
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		rmq.config.User, rmq.config.Password, rmq.config.Host, rmq.config.Port, rmq.config.VHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	rmq.mu.Lock()
	defer rmq.mu.Unlock()

	rmq.conn = conn
	rmq.channel = ch
	rmq.isConnected = true

	for exchange, exchangeType := range rmq.exchanges {
		if err := rmq.declareExchange(exchange, exchangeType); err != nil {
			return err
		}
	}

	if err := rmq.redeclareQueues(); err != nil {
		return err
	}

	if err := rmq.rebindQueues(); err != nil {
		return err
	}

	if err := rmq.channel.QueueBind(
		rmq.constants.Queue.DLQ,
		"",
		rmq.constants.Exchange.DLX,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (rmq *RabbitMQ) redeclareQueues() error {
	for queue := range rmq.queues {
		if queue == rmq.constants.Queue.DLQ {
			if _, err := rmq.channel.QueueDeclare(
				queue,
				true,
				false,
				false,
				false,
				nil,
			); err != nil {
				return err
			}
		} else if queue == rmq.constants.Events.NotificationsEmail || queue == rmq.constants.Events.NotificationsPush || queue == rmq.constants.Events.UserRegistered || queue == rmq.constants.Events.SendNotification {
			if err := rmq.declareQueueWithDLX(queue, rmq.constants.Exchange.DLX); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rmq *RabbitMQ) rebindQueues() error {
	for queue, routingKeys := range rmq.bindings {
		for _, routingKey := range routingKeys {
			if queue == rmq.constants.Events.NotificationsEmail || queue == rmq.constants.Events.NotificationsPush || queue == rmq.constants.Events.UserRegistered || queue == rmq.constants.Events.SendNotification {
				if err := rmq.bindQueue(queue, rmq.constants.Exchange.Notifications, routingKey); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (rmq *RabbitMQ) PublishMessage(queue string, message interface{}) error {
	rmq.mu.RLock()
	connected := rmq.isConnected
	rmq.mu.RUnlock()

	if !connected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx := context.Background()
	err = rmq.channel.PublishWithContext(
		ctx,
		rmq.constants.Exchange.Notifications,
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

func (rmq *RabbitMQ) ConsumeMessages(queue string, handler func([]byte) error) error {
	rmq.mu.RLock()
	connected := rmq.isConnected
	rmq.mu.RUnlock()

	if !connected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	msgs, err := rmq.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error processing message: %v", err)
				rmq.handleRetry(&d, err)
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (rmq *RabbitMQ) handleRetry(d *amqp.Delivery, processingErr error) {
	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}

	retryCount, ok := d.Headers[rmq.constants.Headers.RetryCount].(int)
	if !ok {
		retryCount = 0
	}
	retryCount++
	d.Headers[rmq.constants.Headers.RetryCount] = retryCount
	d.Headers[rmq.constants.Headers.LastError] = processingErr.Error()

	log.Printf("Retrying message, retry count: %d", retryCount)

	if retryCount >= rmq.config.MaxRetryCount {
		log.Printf("Max retries exceeded, rejecting message")
		d.Reject(false)
		return
	}

	delay := time.Duration(retryCount) * time.Second
	time.Sleep(delay)

	err := rmq.channel.PublishWithContext(
		context.Background(),
		d.Exchange,
		d.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Headers:      d.Headers,
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Printf("Failed to republish message: %v", err)
		d.Reject(false)
		return
	}

	d.Ack(false)
}

func (rmq *RabbitMQ) Close() error {
	close(rmq.closeChannel)

	if rmq.channel != nil {
		if err := rmq.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}

	if rmq.conn != nil {
		if err := rmq.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	return nil
}
