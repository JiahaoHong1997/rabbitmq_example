package main

import "rabbitmq_example/RabbitMQ"

func main() {
	hjhOne := RabbitMQ.NewRabbitMQRouting("exOne", "key_two")
	hjhOne.RecieveRouting()
}