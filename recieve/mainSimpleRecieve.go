package main

import "rabbitmq_example/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("HongJiahaoSimple")
	rabbitmq.ConsumeSimple()
}