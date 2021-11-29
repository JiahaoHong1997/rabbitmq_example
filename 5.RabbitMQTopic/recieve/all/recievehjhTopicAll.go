package main

import "rabbitmq_example/RabbitMQ"

func main() {
	hjhOne := RabbitMQ.NewRabbitMQTopic("exTopic", "#")
	hjhOne.RecieveTopic()
}