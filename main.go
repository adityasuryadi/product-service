package main

import (
	"errors"
	"log"
	"product-service/config"
	"product-service/controller"
	"product-service/exception"
	"product-service/pkg/postgres"
	rabbitmq "product-service/pkg/rabbitmq"
	"product-service/repository"
	"product-service/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func ListenQueue() {
	exchange := "order.created"
	queue := "order.create"
	routingKey := "create"

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
	}
	log.Print("producer: declaring binding")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func main() {
	cfg, err := config.LoadConfig("dev.env")
	if err != nil {
		panic("canot load config file")
	}
	app := fiber.New()

	exception.FailOnError(errors.New("INTERNAL_SERVER_ERROR"), "cannot open rabbitmq")

	db := postgres.NewConnPostgres(cfg)
	rabbitMQConn, err := rabbitmq.NewRabbitMqConn(cfg)
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository, rabbitMQConn)
	productController := controller.NewProductController(productService)
	defer rabbitMQConn.Close()
	app.Use(recover.New(recover.ConfigDefault))
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, Range, Authorization",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		// AllowCredentials: false,
	}))
	productController.Route(app)
	app.Listen(":5002")
}
