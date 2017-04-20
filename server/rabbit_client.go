package server

import (
	"github.com/streadway/amqp"
)

func NewRabbitClient(conn string) (*amqp.Connection, error) {
	local, err := amqp.Dial(conn)
	return local, err
}
