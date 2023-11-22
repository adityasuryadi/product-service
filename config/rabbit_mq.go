package config

import (
	"context"
	"encoding/json"
	"log"
	"product-service/exception"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MqChannel struct {
	exchange   string
	queue      string
	routingKey string
}

func NewRabbitMqConn(exchange, queue, routingKey string) *MqChannel {
	return &MqChannel{
		exchange:   exchange,
		queue:      queue,
		routingKey: routingKey,
	}
}

func ConnectionClose() {

}

func (c *MqChannel) ChannelDeclare(data interface{}) {
	conn, err := InitRabbitMQ()
	if err != nil {
		exception.FailOnError(err, "failed to connect Rabbit MQ")
	}

	defer conn.Close()

	ch, err := conn.Channel()
	exception.FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	err = ch.ExchangeDeclare(
		c.exchange, // name
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments)
	)
	exception.FailOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q, err := ch.QueueDeclare(
		c.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	exception.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(q.Name, c.routingKey, c.exchange, false, nil)
	if err != nil {
		exception.FailOnError(err, "Failed to declare a queue")
	}

	// body := bodyFrom(os.Args)
	// body := "Hello"
	body, err := json.Marshal(data)
	if err != nil {
		exception.FailOnError(err, "failed convert body")
	}
	err = ch.PublishWithContext(ctx,
		c.exchange,   // exchange
		c.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	exception.FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func ChannelCLose() {

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func InitRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq_container:5672/")
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		return nil, err
	}
	return conn, nil
}
