package main

import "rabbitmq_example/RabbitMQ"

func main() {
	hjhOne := RabbitMQ.NewRabbitMQTopic("exTopic", "hjh.*.two")
	hjhOne.RecieveTopic()
}