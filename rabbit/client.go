package rabbit

import (
	"fmt"
	"log"
	"time"

	"github.com/daparadoks/go-fuel-api/config"
	"github.com/streadway/amqp"
)

type Client struct {
	connection *amqp.Connection
	quesConfig config.QuesConfig
}

func NewRabbitClient(rabbitConfig config.RabbitConfig, quesConfig config.QuesConfig) *Client {
	c := createConnection(rabbitConfig)
	return &Client{
		connection: c,
		quesConfig: quesConfig,
	}
}

func createConnection(rabbitConfig config.RabbitConfig) *amqp.Connection {
	amqpConfig := amqp.Config{
		Properties: amqp.Table{
			"connection_name": rabbitConfig.ConnectionName,
		},
		Heartbeat: 30 * time.Second,
	}
	connectionUrl := getConnectionUrl(rabbitConfig)
	connection, err := amqp.DialConfig(connectionUrl, amqpConfig)
	if err != nil {
		_ = connection.Close()
		log.Panic("Connection failed to rabbit: %s", err.Error())
	}

	return connection
}

func (c *Client) CloseConnection() {
	c.connection.Close()
}

func getConnectionUrl(config config.RabbitConfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.VirtualHost)
}

func (c *Client) DeclareExchangeQueBindings() {
	channel := c.CreateChannel(0)
	configs := c.getRegisteredQues()
	for _, queConfig := range configs {
		declareExchange(channel, queConfig)
		declareQue(channel, queConfig)
		declareDeadLetterQue(channel, queConfig)
		bindQue(channel, queConfig)
		err := channel.Qos(queConfig.PrefetchCount, 0, false)
		if err != nil {
			log.Panicf("Prefetch couldn't defined: %s", err.Error())
		}
	}
}

func (c *Client) CreateChannel(prefetchCount int) *amqp.Channel {
	channel, err := c.connection.Channel()
	if err != nil {
		channel.Close()
		log.Panicf("Prefetch couldn't defined: %s", err.Error())
	}

	e := channel.Qos(prefetchCount, 0, false)
	if e != nil {
		log.Panicf("Prefetch couldn't defined: %s", err.Error())
	}

	return channel
}

func declareExchange(channel *amqp.Channel, queConfig config.QueConfig) {
	err := channel.ExchangeDeclare(queConfig.Exchange, queConfig.ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Panicf("Exchange declaration error: %s", err.Error())
	}
}

func declareQue(channel *amqp.Channel, queConfig config.QueConfig) {
	deadLetterArgs := getDeadLetterArgs(queConfig.Queue)
	_, err := channel.QueueDeclare(queConfig.Queue, true, false, false, false, deadLetterArgs)
	if err != nil {
		log.Panicf("Queue declaration error: %s", err.Error())
	}
}

func declareDeadLetterQue(channel *amqp.Channel, queConfig config.QueConfig) {
	_, err := channel.QueueDeclare(queConfig.Queue+".deadLetter", true, false, false, false, nil)
	if err != nil {
		log.Panicf("Queue declaration error: %s", err.Error())
	}
}

func bindQue(channel *amqp.Channel, queConfig config.QueConfig) {
	err := channel.QueueBind(queConfig.Queue, queConfig.RoutingKey, queConfig.Exchange, false, nil)
	if err != nil {
		log.Panicf("Queue bind error: %s", err.Error())
	}
}

func getDeadLetterArgs(queName string) amqp.Table {
	return amqp.Table{}
}
