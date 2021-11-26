package main

import (
	"fmt"
	"rabbitmq_example/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("HongJiahaoSimple")
	rabbitmq.PublishSimple("Hello HongJiahao!")
	fmt.Println("发送成功！")

}