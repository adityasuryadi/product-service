package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"product-service/exception"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channell *amqp.Channel
	pCfg     *PublisherConfig
}

type PublisherConfig struct {
	Exchange    string
	QueueName   string
	RoutingKey  string
	ConsumerTag string
}

func NewPublisher(amqpConn *amqp.Connection, pCfg *PublisherConfig) *Publisher {
	ch, err := amqpConn.Channel()
	exception.FailOnError(err, "failed create channel")
	// defer amqpConn.Close()
	return &Publisher{
		channell: ch,
		pCfg:     pCfg,
	}
}

func (p *Publisher) SetupExchangeAndQueue() {
	fmt.Println("declare exchange")
	ch := p.channell
	err := ch.ExchangeDeclare(
		p.pCfg.Exchange, // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments)
	)
	exception.FailOnError(err, "Failed to declare an exchange")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	q, err := ch.QueueDeclare(
		p.pCfg.QueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	exception.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(q.Name, p.pCfg.RoutingKey, p.pCfg.Exchange, false, nil)
	exception.FailOnError(err, "Failed to declare a queue")

}

func (p *Publisher) CloseChannel() {
	if err := p.channell.Close(); err != nil {
		exception.FailOnError(err, "Close Channel")
	}
}

func (p *Publisher) Publish(body []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p.channell.PublishWithContext(ctx, p.pCfg.Exchange, p.pCfg.RoutingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
	log.Printf(" [x] Sent %s", body)
}
