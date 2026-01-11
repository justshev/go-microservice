package broker

import "github.com/rabbitmq/amqp091-go"


func Connect(url string ) (*amqp091.Connection, error) {
	return amqp091.Dial(url)

}