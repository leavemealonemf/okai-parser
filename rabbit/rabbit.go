package rabbit

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Conn() *amqp.Connection {
	usr, _ := os.LookupEnv("RMQ_USER")
	pass, _ := os.LookupEnv("RMQ_PASS")
	connStr := fmt.Sprintf("amqp://%s:%s@localhost:5672/", usr, pass)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %s", err)
	}

	return conn
}
