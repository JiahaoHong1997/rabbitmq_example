package main

import "rabbitmq_example/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("newProduct")
	rabbitmq.RecieveSub()
}
