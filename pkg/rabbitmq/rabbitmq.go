package rabbitmq

import (
	"product-service/exception"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func NewRabbitMqConn(cfg *viper.Viper) (*amqp.Connection, error) {
	// connAddr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
	// 	cfg.GetString("RABBITMQ_USER"),
	// 	cfg.GetString("RABBITMQ_PASSWORD"),
	// 	cfg.GetString("RABBITMQ_HOST"),
	// 	cfg.GetString("RABBITMQ_PORT"))
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq_container:5672/")
	exception.FailOnError(err, "unable connect rabbitmq")
	return conn, nil
}
